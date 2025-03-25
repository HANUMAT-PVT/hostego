package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func MessMenuRoutes(app *fiber.App) {
	api := app.Group("/api/mess-menu")
	api.Patch("/:id", controllers.UpdateMessMenu)
	api.Get("/", controllers.FetchMessMenu)
	api.Post("/create", controllers.CreateMessMenuDate)
}
