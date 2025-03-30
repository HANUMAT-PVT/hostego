package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func DeliveryPartnerWalletRoutes(app *fiber.App) {
	deliveryPartnerWalletRoutes := app.Group("/api/delivery-partner-wallet", middlewares.VerifyUserAuthCookieMiddleware())
	deliveryPartnerWalletRoutes.Get("/withdrawal-requests", controllers.FetchAllDeliveryPartnersTransactions)
	deliveryPartnerWalletRoutes.Get("/:id", controllers.FetchDeliveryPartnerWallet)
	deliveryPartnerWalletRoutes.Get("/transactions/:id", controllers.FetchDeliveryPartnerWalletTransactions)
	deliveryPartnerWalletRoutes.Post("/withdrawal-requests", controllers.CreateWalletWithdrawalRequests)
	deliveryPartnerWalletRoutes.Patch("/withdrawal-requests/:transaction_id/verify", controllers.VerifyDeliveryPartnerWithdrawalRequest)
}
