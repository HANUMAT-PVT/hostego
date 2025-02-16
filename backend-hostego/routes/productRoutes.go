package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func ProductRoutes(app *fiber.App) error {
	productRoutes := app.Group("/api/products")

	productRoutes.Get("/all", controllers.FetchProducts)

	productRoutes.Get("/:id", controllers.FetchProductById)

	productRoutes.Post("/", controllers.CreateNewProduct)

	return nil
}
