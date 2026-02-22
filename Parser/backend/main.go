package main

import (
	"Parser/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// η
	r := gin.Default()

	r.Static("/web", "../web")
	r.StaticFile("/", "../web/index.html")

	routes.SetupAnalyzerRoutes(r)

	r.Run("0.0.0.0:5005")

}
