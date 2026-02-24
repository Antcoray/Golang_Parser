package routes

import (
	"Parser/models"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/scanner"
	"go/token"
	"go/types"
	"io"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SetupAnalyzerRoutes(r *gin.Engine) {
	//r.GET("/result", InitializeAnalyzer)
	r.POST("/upload", RunAnalyzer)
}

type Analyzer struct {
	config          models.Config
	operatorCount   map[string]int
	operandCount    map[string]int
	uniqueOperators int
	uniqueOperands  int
	operatorsTotal  int
	operandsTotal   int
	semanticInfo    *types.Info
	fset            *token.FileSet
	inElseIfChain   bool
}

func (v *Analyzer) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		fmt.Printf("\033[37mVisiting: %T at %s\n", node, v.fset.Position(node.Pos()))
	}

	// Вспомогательная функция для добавления оператора с отладкой
	addOp := func(op string) {
		v.operatorCount[op]++
		if v.fset != nil && node != nil {
			pos := v.fset.Position(node.Pos())
			fmt.Printf("\033[32mOperator %q at line %d (type: %T)\n", op, pos.Line, node)
		}
	}

	addOperand := func(name string) {
		v.operandCount[name]++
		if v.fset != nil && node != nil {
			pos := v.fset.Position(node.Pos())
			fmt.Printf("\033[34mOperand %q at line %d (type: %T)\n", name, pos.Line, node)
		}
	}

	switch n := node.(type) {

	// operands
	case *ast.Ident:
		if n.Name == "_" {
			addOperand("_")
			break
		}
		if obj := v.semanticInfo.Defs[n]; obj != nil {
			switch obj.(type) {
			case *types.Var, *types.Const, *types.Nil, *types.PkgName:
				addOperand(n.Name)
			case *types.Func, *types.Builtin:
				addOp(n.Name)
			}
		} else if obj := v.semanticInfo.Uses[n]; obj != nil {
			switch obj.(type) {
			case *types.Var, *types.Const, *types.Nil, *types.PkgName:
				addOperand(n.Name)
			case *types.Func, *types.Builtin:
				addOp(n.Name)
			}
		}

	// ignore imports
	case *ast.ImportSpec:
		return nil

	// operands: literals
	case *ast.BasicLit:
		addOperand(n.Value)

	// ===== ValueSpec – handle = in var/const declarations =====
	case *ast.ValueSpec:
		if len(n.Values) > 0 {
			addOp("=")
		}
		// Traverse all children to collect operands
		for _, name := range n.Names {
			ast.Walk(v, name)
		}
		if n.Type != nil {
			ast.Walk(v, n.Type)
		}
		for _, val := range n.Values {
			ast.Walk(v, val)
		}
		return nil

	// switch-case
	case *ast.TypeSwitchStmt:
		addOp("switch")
		if assign, ok := n.Assign.(*ast.AssignStmt); ok && assign.Tok == token.DEFINE {
			for _, lhs := range assign.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					addOperand(id.Name) // manually count operand defined in switch
				}
			}
		}
	case *ast.SwitchStmt:
		addOp("switch")

	// select
	case *ast.SelectStmt:
		addOp("select")

	// assignment operators
	case *ast.AssignStmt:
		var op string
		switch n.Tok {
		case token.ASSIGN:
			op = "="
		case token.DEFINE:
			op = ":="
		case token.ADD_ASSIGN:
			op = "+="
		case token.SUB_ASSIGN:
			op = "-="
		case token.MUL_ASSIGN:
			op = "*="
		case token.QUO_ASSIGN:
			op = "/="
		case token.REM_ASSIGN:
			op = "%="
		case token.AND_ASSIGN:
			op = "&="
		case token.OR_ASSIGN:
			op = "|="
		case token.XOR_ASSIGN:
			op = "^="
		case token.SHL_ASSIGN:
			op = "<<="
		case token.SHR_ASSIGN:
			op = ">>="
		case token.AND_NOT_ASSIGN:
			op = "&^="
		}
		if op != "" {
			addOp(op)
		}

	// Control constructs: if (consider else-if as single operator)
	case *ast.IfStmt:
		if !v.inElseIfChain {
			addOp("if")
		}
		// Manually traverse children
		if n.Cond != nil {
			ast.Walk(v, n.Cond)
		}
		if n.Body != nil {
			ast.Walk(v, n.Body)
		}
		if n.Else != nil {
			if elseIf, ok := n.Else.(*ast.IfStmt); ok {
				old := v.inElseIfChain
				v.inElseIfChain = true
				ast.Walk(v, elseIf)
				v.inElseIfChain = old
			} else {
				ast.Walk(v, n.Else)
			}
		}
		return nil

	case *ast.ForStmt:
		addOp("for")
	case *ast.RangeStmt:
		addOp("for")
		addOp("range")

	case *ast.GoStmt:
		addOp("go")

	case *ast.DeferStmt:
		addOp("defer")

	case *ast.ReturnStmt:
		addOp("return")

	case *ast.BranchStmt:
		var op string
		switch n.Tok {
		case token.BREAK:
			op = "break"
		case token.CONTINUE:
			op = "continue"
		case token.GOTO:
			op = "goto"
		case token.FALLTHROUGH:
			op = "fallthrough"
		}
		if op != "" {
			addOp(op)
		}

	case *ast.IncDecStmt:
		if n.Tok == token.INC {
			addOp("++")
		} else if n.Tok == token.DEC {
			addOp("--")
		}

	case *ast.SendStmt:
		addOp("<-")

	case *ast.BinaryExpr:
		addOp(n.Op.String())

	case *ast.UnaryExpr:
		addOp(n.Op.String())
	case *ast.StarExpr:
		addOp("*")
	case *ast.KeyValueExpr:
		addOp(":")
	case *ast.LabeledStmt:
		addOp(":")
	case *ast.CaseClause:
		if n.List == nil { // default:
			addOp("default")
		} else {
			addOp("case")
			addOp(":")
		}
	case *ast.CommClause:
		if n.Comm == nil { // default:
			addOp("default")
		} else {
			addOp("case")
			addOp(":")
		}

	case *ast.CallExpr:
		addOp("()")

	case *ast.IndexExpr:
		addOp("[]")

	case *ast.SliceExpr:
		addOp("[ : ]")

	case *ast.TypeAssertExpr:
		addOp(".(type)")

	case *ast.Ellipsis:
		addOp("...")

	// ========== Halstead operators ==========
	// Dot (selector)
	case *ast.SelectorExpr:
		addOp(".")

	// Block
	case *ast.BlockStmt:
		addOp("{}")

	// Parentheses for grouping
	case *ast.ParenExpr:
		addOp("()")

	case *ast.StructType:
		addOp("struct")
		addOp("{}") // скобки структуры

	case *ast.InterfaceType:
		addOp("interface")
		addOp("{}") // скобки интерфейса

	case *ast.CompositeLit:
		addOp("{}")

	case *ast.ArrayType:
		addOp("[]") // array type also uses []

	case *ast.MapType:
		addOp("map")

	case *ast.ChanType:
		addOp("chan")

	case *ast.FuncType:
		// Parentheses around parameters (always present)
		addOp("()")
		// If results are parenthesized (e.g., multiple return values)
		if n.Results != nil && n.Results.Opening.IsValid() {
			addOp("()")
		}

	case *ast.FuncDecl:
		addOp("func")
		if n.Name != nil {
			ast.Walk(v, n.Name)
		}
		if n.Recv != nil {
			ast.Walk(v, n.Recv)
		}
		if n.Type != nil {
			ast.Walk(v, n.Type)
		}
		if n.Body != nil {
			ast.Walk(v, n.Body)
		}
		return nil

	case *ast.GenDecl:
		var op string
		switch n.Tok {
		case token.VAR:
			op = "var"
		case token.TYPE:
			op = "type"
		case token.CONST:
			op = "const"
		default:
			return v // ignore import and others
		}
		addOp(op)
		// Traverse specifications inside
		for _, spec := range n.Specs {
			ast.Walk(v, spec)
		}
		return nil
	}

	return v
}

func CalculateHalsteadMetrics(v *Analyzer) {
	for _, val := range v.operatorCount {
		v.operatorsTotal += val
		v.uniqueOperators++
	}

	for _, val := range v.operandCount {
		v.operandsTotal += val
		v.uniqueOperands++
	}
}

func RunAnalyzer(c *gin.Context) {
	fpath := filepath.Join("configs", "config.json")
	config := models.LoadConfig(fpath)
	fmt.Println(config)

	analyzer := &Analyzer{
		operatorCount: make(map[string]int),
		operandCount:  make(map[string]int),
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		panic(err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, header.Filename, content, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
		Uses: make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{Importer: importer.Default()}
	_, err = conf.Check(node.Name.Name, fset, []*ast.File{node}, info)
	if err != nil {
		panic(err)
	}

	analyzer.semanticInfo = info
	analyzer.fset = fset

	// Считаем точки с запятой (если нужно — можно удалить, но оставим)
	analyzer.countPunctuation(header.Filename, content)

	ast.Walk(analyzer, node)

	CalculateHalsteadMetrics(analyzer)

	fmt.Println(analyzer.operatorCount)
	fmt.Println(analyzer.operandCount)
	fmt.Println("Unique operators:", analyzer.uniqueOperators)
	fmt.Println("Unique operands:", analyzer.uniqueOperands)
	fmt.Println("Total operators:", analyzer.operatorsTotal)
	fmt.Println("Total operands:", analyzer.operandsTotal)

	c.JSON(200, gin.H{
		"operators":        analyzer.operatorCount,
		"operands":         analyzer.operandCount,
		"unique_operators": analyzer.uniqueOperators,
		"unique_operands":  analyzer.uniqueOperands,
		"operators_total":  analyzer.operatorsTotal,
		"operands_total":   analyzer.operandsTotal,
	})
}

func (v *Analyzer) countPunctuation(filename string, content []byte) {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, -1, len(content))
	s.Init(file, content, nil, 0)

	semicolonCount := 0
	commaCount := 0
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok == token.SEMICOLON && lit == ";" {
			semicolonCount++
		}
		if tok == token.COMMA {
			commaCount++
		}
	}
	if semicolonCount != 0 {
		v.operatorCount[";"] = semicolonCount
		fmt.Printf("Operator %q count: %d\n", ";", semicolonCount)
	}
	if commaCount != 0 {
		v.operatorCount[","] = commaCount
		fmt.Printf("Operator %q count: %d\n", ",", commaCount)
	}
}
