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

func setup(c *gin.Context) *Analyzer {
	fpath := filepath.Join("configs", "config.json")
	config := models.LoadConfig(fpath)
	fmt.Println(config)

	Analyzer := &Analyzer{
		operatorCount:        make(map[string]int),
		operandCount:         make(map[string]int),
		operatorCountClassic: make(map[string]int),
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

	Analyzer.semanticInfo = info
	Analyzer.fset = fset

	Analyzer.countPunctuation(header.Filename, content)

	ast.Walk(Analyzer, node)

	return Analyzer
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
		v.operatorCountClassic[";"] = semicolonCount
		fmt.Printf("Operator %q count: %d\n", ";", semicolonCount)
	}
	if commaCount != 0 {
		v.operatorCount[","] = commaCount
		v.operatorCountClassic[","] = commaCount
		fmt.Printf("Operator %q count: %d\n", ",", commaCount)
	}
}

func (v *Analyzer) CalculateHalsteadMetrics() {
	for _, val := range v.operatorCount {
		v.operatorsTotal += val
		v.uniqueOperators++
	}

	for _, val := range v.operandCount {
		v.operandsTotal += val
		v.uniqueOperands++
	}
}
