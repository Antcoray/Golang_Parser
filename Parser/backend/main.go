package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

type Config struct {
	Operators map[string]bool `json:"operators"`
}

type Analyzer struct {
	config   Config
	countmap map[string]int
}

func (v *Analyzer) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.AssignStmt:
		if n.Tok == token.ASSIGN && v.config.Operators["="] {
			v.countmap["="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.DEFINE && v.config.Operators[":="] {
			v.countmap[":="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.ADD_ASSIGN && v.config.Operators["+="] {
			v.countmap["+="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.SUB_ASSIGN && v.config.Operators["-="] {
			v.countmap["-="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.MUL_ASSIGN && v.config.Operators["*="] {
			v.countmap["*="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.QUO_ASSIGN && v.config.Operators["/="] {
			v.countmap["/="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.REM_ASSIGN && v.config.Operators["%="] {
			v.countmap["%="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.AND_ASSIGN && v.config.Operators["&="] {
			v.countmap["&="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.OR_ASSIGN && v.config.Operators["|="] {
			v.countmap["|="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.XOR_ASSIGN && v.config.Operators["^="] {
			v.countmap["^="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.SHL_ASSIGN && v.config.Operators["<<="] {
			v.countmap["<<="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.SHR_ASSIGN && v.config.Operators[">>="] {
			v.countmap[">>="]++
			fmt.Println(n.Tok)
		}
		if n.Tok == token.AND_NOT_ASSIGN && v.config.Operators["&^="] {
			v.countmap["&^="]++
			fmt.Println(n.Tok)
		}
	}
	return v
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

	analyzer := Analyzer{countmap: make(map[string]int)}

	fpath = filepath.Join("..", "..", "ExampleCode", "main.go")

	node, err := parser.ParseFile(token.NewFileSet(), fpath, nil, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	analyzer.config = config

	ast.Walk(&analyzer, node)

	fmt.Println(analyzer.countmap)

}
