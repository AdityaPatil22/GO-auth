package main

import (
	"GO-temp-backend/config"
	"GO-temp-backend/routes"
	"GO-temp-backend/services"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()
	services.InitServices()

	app := SetupApp()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}

func SetupApp() *fiber.App {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running")
	})
	routes.AuthRoutes(app)

	return app
}
