package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func DeliveryPartnerWalletRoutes(app *fiber.App) {
	deliveryPartnerWalletRoutes := app.Group("/api/delivery-partner-wallet")
	deliveryPartnerWalletRoutes.Get("/withdrawal-requests", controllers.FetchAllDeliveryPartnersTransactions)
	deliveryPartnerWalletRoutes.Get("/:id", controllers.FetchDeliveryPartnerWallet)
	deliveryPartnerWalletRoutes.Get("/transactions/:id", controllers.FetchDeliveryPartnerWalletTransactions)
	deliveryPartnerWalletRoutes.Patch("/withdrawal-requests/:transaction_id/verify", controllers.VerifyDeliveryPartnerWithdrawalRequest)
}
