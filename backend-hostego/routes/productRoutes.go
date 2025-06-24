package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app *fiber.App) error {
	productRoutes := app.Group("/api/products", middlewares.VerifyUserAuthCookieMiddleware())

	productRoutes.Get("/all", controllers.FetchProducts)

	productRoutes.Get("/:id", controllers.FetchProductById)

	productRoutes.Post("/",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin","admin","inventory_manager"),
		controllers.CreateNewProduct,
	)

	productRoutes.Patch("/:id",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin","admin","inventory_manager"),
		controllers.UpdateProductById,
	)
	return nil
}
