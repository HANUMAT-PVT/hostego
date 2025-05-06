package routes

import (
	"backend-hostego/internal/app/hostego-service/controllers"

	"github.com/gofiber/fiber/v3"
)

func ShopRoutes(app *fiber.App) error {

	shopRoutes := app.Group("/api/shop")

	shopRoutes.Post("/", controllers.CreateShop)

	shopRoutes.Get("/:id", controllers.FetchShopById)

	shopRoutes.Get("/", controllers.FetchShops)

	return nil
}
