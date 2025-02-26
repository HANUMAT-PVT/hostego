package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func AddressRoutes(app *fiber.App) {
	addressRoutes := app.Group("/api/address")
	addressRoutes.Get("/", controllers.FetchUserAddress)
	addressRoutes.Post("/", controllers.CreateNewAddress)
	addressRoutes.Patch("/:id", controllers.UpdateAddress)
}
