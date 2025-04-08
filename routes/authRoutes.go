package routes

import (
	"GO-temp-backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/signup", controllers.Signup)
	auth.Post("/login", controllers.Login)
}
