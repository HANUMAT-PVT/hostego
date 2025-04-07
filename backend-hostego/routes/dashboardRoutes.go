package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v3"
)

func DashboardRoutes(app *fiber.App) {
	dashboardRoutes := app.Group("/api/dashboard", middlewares.VerifyUserAuthCookieMiddleware())
	dashboardRoutes.Get("/stats",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin"),
		controllers.GetDashBoardStats)
}
