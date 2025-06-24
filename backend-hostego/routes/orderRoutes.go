package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func OrderRoutes(app *fiber.App) {
	orderRoutes := app.Group("/api/order", middlewares.VerifyUserAuthCookieMiddleware())

	orderRoutes.Get("/order-items",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "admin", "inventory_manager"),
		controllers.FetchAllOrderItemsAccordingToProducts,
	)
	orderRoutes.Post("/cancel-no-refund",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin"),
		controllers.CancelOrder,
	)

	orderRoutes.Get("/all",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "order_manager"),
		controllers.FetchAllOrders,
	)

	orderRoutes.Post("/", controllers.CreateNewOrder)
	orderRoutes.Get("/", controllers.FetchAllUserOrders)
	orderRoutes.Get("/:id", controllers.FetchOrderById)
	orderRoutes.Patch("/:id", controllers.UpdateOrderById)

	orderRoutes.Post("/mark-order-delivered", controllers.MarkOrderAsDelivered)

	orderRoutes.Post("/assign-order-delivery",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "order_assign_manager"),
		controllers.AssignOrderToDeliveryPartner)

	orderRoutes.Get("/delivery-partner/:id", controllers.FetchAllOrdersByDeliveryPartner)
}
