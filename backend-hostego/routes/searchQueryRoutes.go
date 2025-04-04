package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func SearchQueryRoutes(app *fiber.App) {
	searchQuery := app.Group("/api/search-query", middlewares.VerifyUserAuthCookieMiddleware())
	searchQuery.Get("/",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "admin", "inventory_manager"),
		controllers.FetchSearchQuery,
	)
}
