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

	walletRoutes.Post("/verifiy-wallet-transaction/:id",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "payments_manager", "admin"),
		controllers.VerifyWalletTransactionById,
	)
	walletRoutes.Get("/transactions", controllers.FetchUserWalletTransactions)

	walletRoutes.Get("/all-transactions",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "payments_manager"),
		controllers.FetchAllWalletTransactions,
	)
	walletRoutes.Get("/users-wallet-balances",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "payments_manager","admin"),
		controllers.FetchUsersWithPositiveWalletBalance)
}
