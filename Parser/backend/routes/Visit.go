package routes

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
)

func (v *Analyzer) Visit(node ast.Node) ast.Visitor {
	if node != nil {
		fmt.Printf("\033[37mVisiting: %T at %s\n", node, v.fset.Position(node.Pos()))
	}

	addOp := func(op string) {
		v.operatorCount[op]++
		if v.fset != nil && node != nil {
			pos := v.fset.Position(node.Pos())
			fmt.Printf("\033[32mOperator %q at line %d (type: %T)\n", op, pos.Line, node)
		}
	}

	addClassicOp := func(op string) {
		v.operatorCountClassic[op]++
		if v.fset != nil && node != nil {
			pos := v.fset.Position(node.Pos())
			fmt.Printf("\033[33mOperator %q at line %d (type: %T)\n", op, pos.Line, node)
		}
	}

	addOperand := func(name string) {
		v.operandCount[name]++
		if v.fset != nil && node != nil {
			pos := v.fset.Position(node.Pos())
			fmt.Printf("\033[34mOperand %q at line %d (type: %T)\n", name, pos.Line, node)
		}
	}

	updateMaxDepth := func() {
		if v.currentDepth > v.maxDepth {
			v.maxDepth = v.currentDepth
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

	case *ast.ValueSpec:
		if len(n.Values) > 0 {
			addOp("=")
			addClassicOp("=")
		}
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
		//addClassicOp("switch")
		updateMaxDepth()
		if assign, ok := n.Assign.(*ast.AssignStmt); ok && assign.Tok == token.DEFINE {
			for _, lhs := range assign.Lhs {
				if id, ok := lhs.(*ast.Ident); ok {
					addOperand(id.Name) // manually count operand defined in switch
				}
			}
		}

		if n.Body != nil {
			cases := 0
			for _, stmt := range n.Body.List {
				if _, ok := stmt.(*ast.CaseClause); ok {
					cases++
				}
			}
			if cases > 0 {
				v.cl += cases - 1 // n ветвей эквивалентны n-1 if

				for c := cases - 1; c > 0; c-- {
					addClassicOp("if")
				}
				if cases > 2 {
					v.currentDepth += cases - 2
					//updateMaxDepth()
					defer func() {
						v.currentDepth -= cases - 2
					}()
				}
			}
		}

		if n.Assign != nil {
			ast.Walk(v, n.Assign)
		}
		if n.Body != nil {
			ast.Walk(v, n.Body)
		}
		return nil

	case *ast.SwitchStmt:
		addOp("switch")
		//addClassicOp("switch")
		if n.Body != nil {
			cases := 0
			for _, stmt := range n.Body.List {
				if _, ok := stmt.(*ast.CaseClause); ok {
					cases++
				}
			}
			if cases > 0 {
				v.cl += cases - 1 // n ветвей эквивалентны n-1 if

				for c := cases - 1; c > 0; c-- {
					addClassicOp("if")
				}
				if cases > 2 {
					v.currentDepth += cases - 2
					//updateMaxDepth()

					defer func() {
						v.currentDepth -= cases - 2
					}()
				}
			}
		}
		if n.Init != nil {
			ast.Walk(v, n.Init)
		}
		if n.Tag != nil {
			ast.Walk(v, n.Tag)
		}
		if n.Body != nil {
			ast.Walk(v, n.Body)
		}
		return nil

	case *ast.SelectStmt:
		addOp("select")
		//addClassicOp("select")
		updateMaxDepth()
		if n.Body != nil {
			cases := 0
			for _, stmt := range n.Body.List {
				if _, ok := stmt.(*ast.CaseClause); ok {
					cases++
				}
			}
			if cases > 0 {
				v.cl += cases - 1 // n ветвей эквивалентны n-1 if

				for c := cases - 1; c > 0; c-- {
					addClassicOp("if")
				}
				if cases > 2 {
					v.currentDepth += cases - 2
					//updateMaxDepth()

					defer func() {
						v.currentDepth -= cases - 2
					}()
				}
			}
		}
		if n.Body != nil {
			ast.Walk(v, n.Body)
		}
		return nil

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
			addClassicOp(op)
		}

	case *ast.IfStmt:
		v.cl++
		addClassicOp("if")
		updateMaxDepth()
		if !v.inElseIfChain {
			addOp("if")
		}
		if n.Cond != nil {
			ast.Walk(v, n.Cond)
		}
		v.currentDepth++
		if n.Body != nil {
			ast.Walk(v, n.Body)
		}
		v.currentDepth--
		if n.Else != nil {
			if elseIf, ok := n.Else.(*ast.IfStmt); ok {
				old := v.inElseIfChain
				v.inElseIfChain = true
				ast.Walk(v, elseIf)
				v.inElseIfChain = old
			} else {
				updateMaxDepth()
				v.currentDepth++
				ast.Walk(v, n.Else)
				v.currentDepth--
			}
		}
		return nil

	case *ast.ForStmt:
		addOp("for")
		addClassicOp("for")
		updateMaxDepth()

		if n.Init != nil {
			ast.Walk(v, n.Init)
		}
		if n.Cond != nil {
			ast.Walk(v, n.Cond)
		}
		if n.Post != nil {
			ast.Walk(v, n.Post)
		}

		v.currentDepth++

		if n.Body != nil {
			ast.Walk(v, n.Body)
		}

		v.currentDepth--
		return nil
	case *ast.RangeStmt:
		addOp("for")
		addOp("range")
		addClassicOp("for")
		addClassicOp("range")
		updateMaxDepth()
		if n.Key != nil {
			ast.Walk(v, n.Key)
		}
		if n.Value != nil {
			ast.Walk(v, n.Value)
		}
		if n.X != nil {
			ast.Walk(v, n.X)
		}
		switch n.Tok {
		case token.DEFINE:
			addOp(":=")
			addClassicOp(":=")
		case token.ASSIGN:
			addOp("=")
			addClassicOp("=")
		}
		v.currentDepth++
		if n.Body != nil {
			ast.Walk(v, n.Body)
		}
		v.currentDepth--
		return nil

	case *ast.GoStmt:
		addOp("go")
		addClassicOp("go")

	case *ast.DeferStmt:
		addOp("defer")
		addClassicOp("defer")

	case *ast.ReturnStmt:
		addOp("return")
		addClassicOp("return")

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
			addClassicOp(op)
		}

	case *ast.IncDecStmt:
		switch n.Tok {
		case token.INC:
			addOp("++")
			addClassicOp("++")
		case token.DEC:
			addOp("--")
			addClassicOp("--")
		}

	case *ast.SendStmt:
		addOp("<-")
		addClassicOp("<-")

	case *ast.BinaryExpr:
		addOp(n.Op.String())
		addClassicOp(n.Op.String())

	case *ast.UnaryExpr:
		addOp(n.Op.String())
		addClassicOp(n.Op.String())
	case *ast.StarExpr:
		addOp("*")
		addClassicOp("*")
	case *ast.KeyValueExpr:
		addOp(":")
		addClassicOp(":")
	case *ast.LabeledStmt:
		addOp(":")
		addClassicOp(":")

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
		addClassicOp("()")

	case *ast.IndexExpr:
		addOp("[]")
		addClassicOp("[]")

	case *ast.SliceExpr:
		addOp("[ : ]")
		addClassicOp("[ : ]")

	case *ast.TypeAssertExpr:
		addOp(".(type)")
		addClassicOp(".(type)")

	case *ast.Ellipsis:
		addOp("...")
		addClassicOp("...")

	case *ast.SelectorExpr:
		addOp(".")
		addClassicOp(".")

	case *ast.BlockStmt:
		addOp("{}")
		addClassicOp("{}")

	case *ast.ParenExpr:
		addOp("()")
		addClassicOp("()")

	case *ast.StructType:
		addOp("struct")
		addOp("{}")
		addClassicOp("struct")
		addClassicOp("{}")

	case *ast.InterfaceType:
		addOp("interface")
		addOp("{}")
		addClassicOp("interface")
		addClassicOp("{}")

	case *ast.CompositeLit:
		addOp("{}")
		addClassicOp("{}")

	case *ast.ArrayType:
		addOp("[]")
		addClassicOp("[]")

	case *ast.MapType:
		addOp("map")
		addClassicOp("map")

	case *ast.ChanType:
		addOp("chan")
		addClassicOp("chan")

	case *ast.FuncType:
		addOp("()")
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
			return v
		}
		addOp(op)
		addClassicOp(op)
		for _, spec := range n.Specs {
			ast.Walk(v, spec)
		}
		return nil
	}

	return v
}
