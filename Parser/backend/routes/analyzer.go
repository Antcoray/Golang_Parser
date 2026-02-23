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
	switch n := node.(type) {

	//operands
	case *ast.Ident:
		if n.Name == "_" {
			v.operandCount["_"]++
			break
		}

		if obj := v.semanticInfo.Defs[n]; obj != nil {
			switch obj.(type) {
			case *types.Var, *types.Const, *types.Nil:
				v.operandCount[n.Name]++
			case *types.Func:
				v.operatorCount[n.Name]++
			}
		} else if obj := v.semanticInfo.Uses[n]; obj != nil {
			switch obj.(type) {
			case *types.Var, *types.Const, *types.Nil:
				v.operandCount[n.Name]++
			case *types.Func:
				v.operatorCount[n.Name]++
			}
		}

	//ignore imports
	case *ast.ImportSpec:
		return nil

	// operands: literals
	case *ast.BasicLit:
		v.operandCount[n.Value]++

	//switch-case
	case *ast.TypeSwitchStmt:
		if v.config.Operators["switch"] {
			v.operatorCount["switch"]++
		}
		if assign, ok := n.Assign.(*ast.AssignStmt); ok && assign.Tok == token.DEFINE {
			for _, lhs := range assign.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					v.operandCount[id.Name]++ // manually count operand defined in switch
				}
			}
		}
	case *ast.SwitchStmt:
		if v.config.Operators["switch"] {
			v.operatorCount["switch"]++
		}

	// select
	case *ast.SelectStmt:
		if v.config.Operators["select"] {
			v.operatorCount["select"]++
		}

	//assignment operators
	case *ast.AssignStmt:
		if n.Tok == token.ASSIGN && v.config.Operators["="] {
			v.operatorCount["="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.DEFINE && v.config.Operators[":="] {
			v.operatorCount[":="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.ADD_ASSIGN && v.config.Operators["+="] {
			v.operatorCount["+="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.SUB_ASSIGN && v.config.Operators["-="] {
			v.operatorCount["-="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.MUL_ASSIGN && v.config.Operators["*="] {
			v.operatorCount["*="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.QUO_ASSIGN && v.config.Operators["/="] {
			v.operatorCount["/="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.REM_ASSIGN && v.config.Operators["%="] {
			v.operatorCount["%="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.AND_ASSIGN && v.config.Operators["&="] {
			v.operatorCount["&="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.OR_ASSIGN && v.config.Operators["|="] {
			v.operatorCount["|="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.XOR_ASSIGN && v.config.Operators["^="] {
			v.operatorCount["^="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.SHL_ASSIGN && v.config.Operators["<<="] {
			v.operatorCount["<<="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.SHR_ASSIGN && v.config.Operators[">>="] {
			v.operatorCount[">>="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.AND_NOT_ASSIGN && v.config.Operators["&^="] {
			v.operatorCount["&^="]++
			fmt.Println(n.Tok)
		}
		//
		//
		// Управляющие конструкции: if (с учётом else-if как одного оператора)
	case *ast.IfStmt:
		if !v.inElseIfChain && v.config.Operators["if"] {
			v.operatorCount["if"]++
		}
		// Обходим потомков вручную
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

	case *ast.ForStmt, *ast.RangeStmt:
		if v.config.Operators["for"] {
			v.operatorCount["for"]++
		}

	case *ast.GoStmt:
		if v.config.Operators["go"] {
			v.operatorCount["go"]++
		}

	case *ast.DeferStmt:
		if v.config.Operators["defer"] {
			v.operatorCount["defer"]++
		}

	case *ast.ReturnStmt:
		if v.config.Operators["return"] {
			v.operatorCount["return"]++
		}

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
		if op != "" && v.config.Operators[op] {
			v.operatorCount[op]++
		}

	case *ast.IncDecStmt:
		if n.Tok == token.INC && v.config.Operators["++"] {
			v.operatorCount["++"]++
		} else if n.Tok == token.DEC && v.config.Operators["--"] {
			v.operatorCount["--"]++
		}

	case *ast.SendStmt:
		if v.config.Operators["<-"] {
			v.operatorCount["<-"]++
		}

	case *ast.BinaryExpr:
		op := n.Op.String()
		if v.config.Operators[op] {
			v.operatorCount[op]++
		}

	case *ast.UnaryExpr:
		op := n.Op.String()
		if v.config.Operators[op] {
			v.operatorCount[op]++
		}

	case *ast.CallExpr:
		if v.config.Operators["()"] {
			v.operatorCount["()"]++
		}

	case *ast.IndexExpr:
		if v.config.Operators["[]"] {
			v.operatorCount["[]"]++
		}

	case *ast.SliceExpr:
		if v.config.Operators["[ : ]"] {
			v.operatorCount["[ : ]"]++
		}

	case *ast.TypeAssertExpr:
		if v.config.Operators[".(type)"] {
			v.operatorCount[".(type)"]++
		}

	case *ast.Ellipsis:
		if v.config.Operators["..."] {
			v.operatorCount["..."]++
		}

	// ========== Операторы из определения Холстеда ==========
	// Точка (селектор)
	case *ast.SelectorExpr:
		if v.config.Operators["."] {
			v.operatorCount["."]++
		}

	// Составной оператор (блок {})
	case *ast.BlockStmt:
		if v.config.Operators["{}"] {
			v.operatorCount["{}"]++
		}

	// Круглые скобки для группировки
	case *ast.ParenExpr:
		if v.config.Operators["()"] {
			v.operatorCount["()"]++

		}

	case *ast.StructType:
		if v.config.Operators["struct"] {
			v.operatorCount["struct"]++
		}
		// поля обрабатываются автоматически

	case *ast.InterfaceType:
		if v.config.Operators["interface"] {
			v.operatorCount["interface"]++
		}

	case *ast.ArrayType:
		if v.config.Operators["[]"] {
			v.operatorCount["[]"]++ // тип массива тоже использует []
		}

	case *ast.MapType:
		if v.config.Operators["map"] {
			v.operatorCount["map"]++
		}

	case *ast.ChanType:
		if v.config.Operators["chan"] {
			v.operatorCount["chan"]++
		}

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

	analyzer := Analyzer{operatorCount: make(map[string]int), operandCount: make(map[string]int)}

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
	//
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
	//
	analyzer.config = config

	analyzer.countSemicolons(header.Filename, content)

	ast.Walk(&analyzer, node)

	CalculateHalsteadMetrics(&analyzer)

	fmt.Println(analyzer.operatorCount)
	fmt.Println(analyzer.operandCount)
	fmt.Println(analyzer.uniqueOperators)
	fmt.Println(analyzer.uniqueOperands)
	fmt.Println(analyzer.operatorsTotal)
	fmt.Println(analyzer.operandsTotal)

	c.JSON(200, gin.H{
		"operators":        analyzer.operatorCount,
		"operands":         analyzer.operandCount,
		"unique_operators": analyzer.uniqueOperators,
		"unique_operands":  analyzer.uniqueOperands,
		"operators_total":  analyzer.operatorsTotal,
		"operands_total":   analyzer.operandsTotal,
	})
}

func (v *Analyzer) countSemicolons(filename string, content []byte) {
	if !v.config.Operators[";"] {
		return
	}

	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile(filename, -1, len(content))
	s.Init(file, content, nil, 0)

	count := 0
	for {
		_, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}

		if tok == token.SEMICOLON && lit == ";" {
			count++
		}
	}
	if count != 0 {
		v.operatorCount[";"] = count
	}
}
