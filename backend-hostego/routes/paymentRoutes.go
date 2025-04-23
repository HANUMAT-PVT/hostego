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
		middlewares.RoleMiddleware("super_admin", "payments_manager", "admin", "order_manager"),
		controllers.InitiateRefundPayment,
	)
	paymentRoutes.Post("/razorpay", controllers.InitateRazorpayPaymentOrder)
	paymentRoutes.Post("/razorpay/verify-payment", controllers.VerifyRazorpayPayment)

	paymentRoutes.Post("/cashfree", controllers.InitateCashfreePaymentOrder)


	paymentRoutes.Post("/cashfree/verify-payment", controllers.VerifyCashfreePayment)
}


func PaymentWebhookRoutes(app * fiber.App){
	webhookRoutes := app.Group("/api/razorpay/webhook")
	// https://backend.hostego.in/api/razorpay/webhook/payment/verify
	webhookRoutes.Post("/payment/verify",controllers.RazorpayWebhookHandler);
}