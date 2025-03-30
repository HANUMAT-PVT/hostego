package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func OrderRoutes(app *fiber.App) {
	orderRoutes := app.Group("/api/order", middlewares.VerifyUserAuthCookieMiddleware())
	orderRoutes.Get("/order-items", controllers.FetchAllOrderItemsAccordingToProducts)
	orderRoutes.Get("/all", controllers.FetchAllOrders)
	orderRoutes.Post("/", controllers.CreateNewOrder)
	orderRoutes.Get("/", controllers.FetchAllUserOrders)
	orderRoutes.Get("/:id", controllers.FetchOrderById)
	orderRoutes.Patch("/:id", controllers.UpdateOrderById)
	orderRoutes.Post("/mark-order-delivered", controllers.MarkOrderAsDelivered)
	orderRoutes.Post("/assign-order-delivery", controllers.AssignOrderToDeliveryPartner)
	orderRoutes.Get("/delivery-partner/:id", controllers.FetchAllOrdersByDeliveryPartner)
}
