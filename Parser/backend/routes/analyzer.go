package routes

import (
	"Parser/models"
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func SetupAnalyzerRoutes(r *gin.Engine) {
	//r.GET("/result", InitializeAnalyzer)
	r.POST("/upload", InitializeAnalyzer)
}

type Analyzer struct {
	config          models.Config
	operatorCount   map[string]int
	operandCount    map[string]int
	uniqueOperators int
	uniqueOperands  int
	operatorsTotal  int
	operandsTotal   int
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

func InitializeAnalyzer(c *gin.Context) {
	fpath := filepath.Join("configs", "config.json")
	config := models.LoadConfig(fpath)
	fmt.Println(config)

	analyzer := Analyzer{operatorCount: make(map[string]int), operandCount: make(map[string]int)}

	fpath = filepath.Join("..", "..", "ExampleCode", "main.go")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		panic(err)
	}

	content, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}

	node, err := parser.ParseFile(token.NewFileSet(), header.Filename, content, parser.SkipObjectResolution)
	if err != nil {
		panic(err)
	}

	analyzer.config = config

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
