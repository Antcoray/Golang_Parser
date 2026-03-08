package routes

import (
	"Parser/models"
	"fmt"
	"go/token"
	"go/types"

	"github.com/gin-gonic/gin"
)

type Analyzer struct {
	config               models.Config
	operatorCount        map[string]int
	operatorCountClassic map[string]int
	operandCount         map[string]int
	uniqueOperators      int
	uniqueOperands       int
	operatorsTotal       int
	operandsTotal        int
	semanticInfo         *types.Info
	fset                 *token.FileSet
	inElseIfChain        bool
	cl                   int
	currentDepth         int
	maxDepth             int
	Gilb                 bool
}

func SetupAnalyzerRoutes(r *gin.Engine) {
	//r.GET("/result", InitializeAnalyzer)
	r.POST("/uploadvol", RunvolAnalyzer)
	r.POST("/uploadflow", RunflowAnalyzer)
}

func RunvolAnalyzer(c *gin.Context) {
	Analyzer := setup(c)

	Analyzer.CalculateHalsteadMetrics()

	fmt.Println(Analyzer.operatorCount)
	fmt.Println(Analyzer.operandCount)
	fmt.Println("Unique operators:", Analyzer.uniqueOperators)
	fmt.Println("Unique operands:", Analyzer.uniqueOperands)
	fmt.Println("Total operators:", Analyzer.operatorsTotal)
	fmt.Println("Total operands:", Analyzer.operandsTotal)

	c.JSON(200, gin.H{
		"operators":        Analyzer.operatorCount,
		"operands":         Analyzer.operandCount,
		"unique_operators": Analyzer.uniqueOperators,
		"unique_operands":  Analyzer.uniqueOperands,
		"operators_total":  Analyzer.operatorsTotal,
		"operands_total":   Analyzer.operandsTotal,
	})
}

func RunflowAnalyzer(c *gin.Context) {
	Analyzer := setup(c)

	ClassicOperatorsCount := 0

	for _, val := range Analyzer.operatorCountClassic {
		ClassicOperatorsCount += val
	}
	fmt.Println(Analyzer.cl)
	fmt.Println(Analyzer.operatorCountClassic)
	fmt.Println(ClassicOperatorsCount)
	c.JSON(200, gin.H{
		"operators": Analyzer.operatorCountClassic,
		"maxDepth":  Analyzer.maxDepth,
		"CL":        Analyzer.cl,
		"cl":        float64(Analyzer.cl) / float64(ClassicOperatorsCount),
	})
}
