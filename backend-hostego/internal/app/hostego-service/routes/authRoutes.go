package routes

import (
	"backend-hostego/internal/app/hostego-service/controllers"

	"github.com/gofiber/fiber/v3"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/api/auth")
	auth.Post("/signup", controllers.Signup)
}
