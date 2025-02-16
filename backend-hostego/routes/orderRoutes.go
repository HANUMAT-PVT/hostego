package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func OrderRoutes(app *fiber.App) {
	orderRoutes := app.Group("/api/order")
	orderRoutes.Post("/", controllers.CreateNewOrder)
	orderRoutes.Get("/:id", controllers.FetchOrderById)
	orderRoutes.Post("/mark-order-delivered", controllers.MarkOrderAsDelivered)
}
