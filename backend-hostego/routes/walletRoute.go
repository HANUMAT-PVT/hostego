package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func WalletRoutes(app *fiber.App) {
	walletRoutes := app.Group("/api/wallet");
	walletRoutes.Get("/", controllers.FetchUserWallet);
	walletRoutes.Post("/credit", controllers.CreditWalletTransaction);
	walletRoutes.Post("/verifiy-wallet-transaction/:id", controllers.VerifyWalletTransactionById);
}
