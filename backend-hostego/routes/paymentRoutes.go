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
	paymentRoutes.Post("/refund", controllers.InitiateRefundPayment)
	paymentRoutes.Post("/cashfree", controllers.InitateCashfreePaymentOrder)
	paymentRoutes.Post("/cashfree/:cf_order_id", controllers.VerifyCashfreePayment)
}
