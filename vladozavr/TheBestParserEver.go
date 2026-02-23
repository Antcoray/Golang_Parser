package main

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"math"
)

// func main() {
func notmain() {
	var operators = map[string]int{
		"+": 0, "-": 0, "*": 0, "/": 0,
		"<=": 0, ">=": 0, "!=": 0, "==": 0,
		">": 0, "<": 0, "=": 0, ":=": 0,
		"<<": 0, ">>": 0, "&": 0, "|": 0,
		"^": 0, "if": 0, "for": 0, "break": 0,
		"continue": 0, "fallthrough": 0, "%": 0, "switch": 0,
		"()": 0,
	}
	var operands = map[string]int{}
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "Test.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println("Error")
	}

	info := types.Info{
		Defs: make(map[*ast.Ident]types.Object),
		Uses: make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{Importer: importer.Default()}
	_, err = conf.Check("mypackage", fset, []*ast.File{node}, &info)
	if err != nil {
		fmt.Println("Error")
	}

	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.BranchStmt:
			{
				switch x.Tok {
				case token.BREAK:
					{
						operators["break"]++
					}
				case token.CONTINUE:
					{
						operators["continue"]++
					}
				case token.FALLTHROUGH:
					{
						operators["fallthrough"]++
					}
				case token.GOTO:
					{
						_, ok := operators["Goto<"+x.Label.Name+">"]
						if !ok {
							operators["Goto<"+x.Label.Name+">"] = 1
						} else {
							operators["Goto<"+x.Label.Name+">"]++
						}
					}
				}
			}
		case *ast.SwitchStmt:
			{
				operators["switch"]++
			}
		case *ast.AssignStmt:
			{
				operators[x.Tok.String()]++
			}
		case *ast.ImportSpec:
			{
				return false
			}
		case *ast.BinaryExpr:
			{
				operators[x.Op.String()]++
			}
		case *ast.IfStmt:
			{
				operators["if"]++
			}
		case *ast.ForStmt:
			{
				operators["for"]++
			}
		}
		return true
	})
	var SlovarProgrammy int = 0
	var DlinaProgrammy int = 0
	var ObemProgrammy int = 0
	for id, obj := range info.Defs {
		_, ok := obj.(*types.Label)
		_, ok1 := obj.(*types.Func)
		if obj != nil && !ok && !ok1 {
			operands[id.Name] = 0
			SlovarProgrammy++
		}
	}
	for id, obj := range info.Uses {
		if obj != nil {
			_, ok := operands[id.Name]
			if ok {
				operands[id.Name]++
			}
		}
		_, ok := obj.(*types.Func)
		if ok {
			_, ok := operators[id.Name]
			if ok {
				operators[id.Name]++
			} else {
				operators[id.Name] = 1
				SlovarProgrammy++
			}
		}
	}
	fmt.Println()
	for key, value := range operands {
		fmt.Println(key, value)
		DlinaProgrammy += value
	}
	fmt.Println()
	fmt.Println()
	for key, value := range operators {
		fmt.Println(key, value)
		if value > 0 {
			SlovarProgrammy++
		}
		DlinaProgrammy += value
	}
	ObemProgrammy = DlinaProgrammy * int(math.Log2(float64(SlovarProgrammy)))
	fmt.Println(SlovarProgrammy, DlinaProgrammy, ObemProgrammy)
}
