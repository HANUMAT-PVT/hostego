package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func PaymentRoutes(app *fiber.App) {
	paymentRoutes := app.Group("/api/payment", middlewares.VerifyUserAuthCookieMiddleware())
	paymentRoutes.Post("/", controllers.InitiatePayment)
	paymentRoutes.Get("/transactions", controllers.FetchUserPaymentTransactions)

	paymentRoutes.Post("/refund",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "payments_manager","admin","order_manager"),
		controllers.InitiateRefundPayment,
	)
	paymentRoutes.Post("/cashfree", controllers.InitateCashfreePaymentOrder)
	paymentRoutes.Post("/cashfree/verify-payment", controllers.VerifyCashfreePayment)
}
