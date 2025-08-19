package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func OrderItemRoutes(app *fiber.App) {
	orderItemRoute := app.Group("/api/order-item", middlewares.VerifyUserAuthCookieMiddleware())

	orderItemRoute.Post("/refund",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "payments_manager"),
		controllers.CancelOrderItemAndInitiateRefund)

	orderItemRoute.Get("/all",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "payments_manager", "admin"),
		controllers.FetchOrderItems)

}

// SELECT
//   u.user_id AS user_id,
//   w.balance,
//   (
//     SELECT SUM(balance)
//     FROM wallets
//     WHERE balance > 0
//   ) AS total_wallet_balance
// FROM
//   users u
// JOIN
//   wallets w ON u.user_id = w.user_id
// WHERE
//   w.balance > 0;
