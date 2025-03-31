package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v3"
)

func DeliveryPartnerWalletRoutes(app *fiber.App) {
	deliveryPartnerWalletRoutes := app.Group("/api/delivery-partner-wallet", middlewares.VerifyUserAuthCookieMiddleware())
	deliveryPartnerWalletRoutes.Get("/withdrawal-requests",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "payments_manager"),
		controllers.FetchAllDeliveryPartnersTransactions,
	)
	deliveryPartnerWalletRoutes.Get("/:id", controllers.FetchDeliveryPartnerWallet)
	deliveryPartnerWalletRoutes.Get("/transactions/:id", controllers.FetchDeliveryPartnerWalletTransactions)

	deliveryPartnerWalletRoutes.Post("/withdrawal-requests",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "payments_manager"),
		controllers.CreateWalletWithdrawalRequests,
	)
	deliveryPartnerWalletRoutes.Patch("/withdrawal-requests/:transaction_id/verify",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "payments_manager"),
		controllers.VerifyDeliveryPartnerWithdrawalRequest,
	)
}
