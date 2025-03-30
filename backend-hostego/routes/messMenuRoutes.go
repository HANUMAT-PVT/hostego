package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func MessMenuRoutes(app *fiber.App) {
	api := app.Group("/api/mess-menu", middlewares.VerifyUserAuthCookieMiddleware())
	api.Patch("/:id", controllers.UpdateMessMenu)
	api.Get("/", controllers.FetchMessMenu)
	api.Post("/create", controllers.CreateMessMenuDate)
}
