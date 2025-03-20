package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)	

func DeliveryPartnerWalletRoutes(app *fiber.App) {
	deliveryPartnerWalletRoutes := app.Group("/delivery-partner-wallet")
	deliveryPartnerWalletRoutes.Get("/", controllers.FetchDeliveryPartnerWallet)
	deliveryPartnerWalletRoutes.Get("/transactions", controllers.FetchDeliveryPartnerWalletTransactions)
	deliveryPartnerWalletRoutes.Post("/withdrawal-requests", controllers.CreateWalletWithdrawalRequests)
	deliveryPartnerWalletRoutes.Post("/withdrawal-requests/:transaction_id/verify", controllers.VerifyDeliveryPartnerWithdrawalRequest)
}
