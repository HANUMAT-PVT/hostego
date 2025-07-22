package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RatingRoutes(app *fiber.App) {
	ratingRoutes := app.Group("api/ratings", middlewares.VerifyUserAuthCookieMiddleware())
	ratingRoutes.Post("/", controllers.SaveRatingAndUpdateStats)
	ratingRoutes.Get("/:product_id", controllers.FetchRatingsForProduct)
	ratingRoutes.Get("/order/:order_id", controllers.FetchRatingsForOrder)
}
