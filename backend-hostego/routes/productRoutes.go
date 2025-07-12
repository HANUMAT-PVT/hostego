package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App) error {
	productRoutes := app.Group("/api/products", middlewares.VerifyUserAuthCookieMiddleware())

	productRoutes.Get("/shop/:shop_id", middlewares.VerifyUserAuthCookieMiddleware(),

		controllers.FetchProductsByShopId,
	)
	productRoutes.Get("/all", controllers.FetchProducts)

	productRoutes.Get("/:id", controllers.FetchProductById)

	productRoutes.Post("/",
		middlewares.VerifyUserAuthCookieMiddleware(),

		controllers.CreateNewProduct,
	)

	productRoutes.Patch("/:id",
		middlewares.VerifyUserAuthCookieMiddleware(),

		controllers.UpdateProductById,
	)

	return nil
}
