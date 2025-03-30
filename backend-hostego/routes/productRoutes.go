package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func ProductRoutes(app *fiber.App) error {
	productRoutes := app.Group("/api/products", middlewares.VerifyUserAuthCookieMiddleware())

	productRoutes.Get("/all", controllers.FetchProducts)

	productRoutes.Get("/:id", controllers.FetchProductById)

	productRoutes.Post("/", controllers.CreateNewProduct)

	productRoutes.Patch("/:id", controllers.UpdateProductById)

	return nil
}
