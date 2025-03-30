package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api", middlewares.VerifyUserAuthCookieMiddleware())

	api.Get("/users", controllers.GetUsers)

	// api.Post("/users", controllers.CreateUser)
}
