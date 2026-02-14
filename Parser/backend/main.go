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
	count int
}

func (v *Analyzer) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.Ident:
		fmt.Println(n.Name)
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

func main() {
	fpath := filepath.Join("Config.json")
	config := LoadConfig(fpath)
	fmt.Println(config)

	analyzer := Analyzer{}

	fpath = filepath.Join("..", "..", "ExampleCode", "main.go")

	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, fpath, nil, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	ast.Walk(&analyzer, node)

}
