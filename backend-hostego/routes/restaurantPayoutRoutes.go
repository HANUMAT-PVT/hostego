package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RestaurantPayoutRoutes(app *fiber.App) {
	restaurantPayoutRoutes := app.Group("/api/restaurant-payout", middlewares.VerifyUserAuthCookieMiddleware())
	restaurantPayoutRoutes.Post("/initiate", controllers.InitiateRestaurantPayout)
	restaurantPayoutRoutes.Patch("/verify/:payout_id", controllers.VerifyRestaurantPayout)
	restaurantPayoutRoutes.Get("/:shop_id", controllers.FetchRestaurantPayouts)
	restaurantPayoutRoutes.Get("/", controllers.GetRestaurantPayouts)
}
