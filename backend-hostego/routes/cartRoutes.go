package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v2"
)

func CartRoutes(app *fiber.App) {
	cartRoutes := app.Group("/api/cart", middlewares.VerifyUserAuthCookieMiddleware())
	cartRoutes.Post("/", controllers.AddProductInUserCart)
	cartRoutes.Get("/", controllers.FetchUserCart)
	cartRoutes.Patch("/:id", controllers.UpdateProductInUserCart)
}
