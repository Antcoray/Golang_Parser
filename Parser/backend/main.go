package main

import (
	"Parser/assets"
	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {

	assets.InitializeAnalyzer()
	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	// Start the server on port 3000
	log.Fatal(app.Listen("0.0.0.0:5005"))

}
