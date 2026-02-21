package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

type Config struct {
	Operators map[string]bool `json:"operators"`
}

type Analyzer struct {
	config        Config
	operatorCount map[string]int
	operandCount  map[string]int
}

func (v *Analyzer) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	//assignment operators
	case *ast.AssignStmt:
		if n.Tok == token.ASSIGN && v.config.Operators["="] {
			v.operatorCount["="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.DEFINE && v.config.Operators[":="] {
			v.operatorCount[":="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.ADD_ASSIGN && v.config.Operators["+="] {
			v.operatorCount["+="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.SUB_ASSIGN && v.config.Operators["-="] {
			v.operatorCount["-="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.MUL_ASSIGN && v.config.Operators["*="] {
			v.operatorCount["*="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.QUO_ASSIGN && v.config.Operators["/="] {
			v.operatorCount["/="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.REM_ASSIGN && v.config.Operators["%="] {
			v.operatorCount["%="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.AND_ASSIGN && v.config.Operators["&="] {
			v.operatorCount["&="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.OR_ASSIGN && v.config.Operators["|="] {
			v.operatorCount["|="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.XOR_ASSIGN && v.config.Operators["^="] {
			v.operatorCount["^="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.SHL_ASSIGN && v.config.Operators["<<="] {
			v.operatorCount["<<="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.SHR_ASSIGN && v.config.Operators[">>="] {
			v.operatorCount[">>="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		if n.Tok == token.AND_NOT_ASSIGN && v.config.Operators["&^="] {
			v.operatorCount["&^="]++
			fmt.Println(n.Tok)
			CountOperands(v, n.Lhs, n.Rhs)
		}
		// other operators
	}
	return v
}

func CountOperands(v *Analyzer, left []ast.Expr, right []ast.Expr) {
	var buff bytes.Buffer
	fset := token.NewFileSet()

	for _, e := range left {
		buff.Reset()
		format.Node(&buff, fset, e)
		v.operandCount[buff.String()]++
	}
	for _, e := range right {
		buff.Reset()
		format.Node(&buff, fset, e)
		v.operandCount[buff.String()]++
	}
}

func LoadConfig(fpath string) Config {

	file, err := os.Open(fpath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}

	return config
}

func (v Analyzer) ToJson() {

}

func main() {
	fpath := filepath.Join("Config.json")
	config := LoadConfig(fpath)
	fmt.Println(config)

	analyzer := Analyzer{operatorCount: make(map[string]int), operandCount: make(map[string]int)}

	fpath = filepath.Join("..", "..", "ExampleCode", "main.go")

	node, err := parser.ParseFile(token.NewFileSet(), fpath, nil, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	analyzer.config = config

	ast.Walk(&analyzer, node)

	fmt.Println(analyzer.operatorCount)
	fmt.Println(analyzer.operandCount)

}
