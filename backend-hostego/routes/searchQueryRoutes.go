package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func SearchQueryRoutes(app *fiber.App) {
	searchQuery := app.Group("/api/search-query")
	searchQuery.Get("/", controllers.FetchSearchQuery)
}
