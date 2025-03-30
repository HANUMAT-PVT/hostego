package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func WalletRoutes(app *fiber.App) {
	walletRoutes := app.Group("/api/wallet", middlewares.VerifyUserAuthCookieMiddleware())
	walletRoutes.Get("/", controllers.FetchUserWallet)
	walletRoutes.Post("/credit", controllers.CreditWalletTransaction)
	walletRoutes.Post("/verifiy-wallet-transaction/:id", controllers.VerifyWalletTransactionById)
	walletRoutes.Get("/transactions", controllers.FetchUserWalletTransactions)
	walletRoutes.Get("/all-transactions", controllers.FetchAllWalletTransactions)
}
