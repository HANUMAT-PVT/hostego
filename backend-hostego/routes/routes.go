package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api")

	api.Get("/users", controllers.GetUsers)

	api.Get("/user/me", controllers.GetUserById)

	api.Post("/users", controllers.CreateUser)

}
