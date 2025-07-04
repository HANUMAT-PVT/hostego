package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func ProductCategoryRoutes(app *fiber.App) {
	productCategoryRoutes := app.Group("/api/product-category", middlewares.VerifyUserAuthCookieMiddleware())

	productCategoryRoutes.Get("/:shop_id", controllers.FetchCategoriesByShopId)
	productCategoryRoutes.Post("/", controllers.CreateNewCategory)

}
