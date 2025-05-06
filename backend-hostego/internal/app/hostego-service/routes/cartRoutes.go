package routes

import (
	"backend-hostego/internal/app/hostego-service/controllers"

	"github.com/gofiber/fiber/v3"
)

func CartRoutes(app *fiber.App) {
	cartRoutes := app.Group("/api/cart")
	// cartRoutes.Post("/", controllers.AddProductInUserCart)
	cartRoutes.Get("/", controllers.FetchUserCart)
	// cartRoutes.Patch("/:id", controllers.UpdateProductInUserCart)
}
