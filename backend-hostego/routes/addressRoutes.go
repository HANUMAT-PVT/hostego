package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AddressRoutes(app *fiber.App) {
	addressRoutes := app.Group("/api/address", middlewares.VerifyUserAuthCookieMiddleware())
	addressRoutes.Get("/", controllers.FetchUserAddress)
	addressRoutes.Post("/", controllers.CreateNewAddress)
	addressRoutes.Patch("/:id", controllers.UpdateAddress)
	addressRoutes.Delete("/:id", controllers.DeleteAddress)
}
