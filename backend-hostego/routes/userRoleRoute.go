package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func UserRoleRoute(app *fiber.App) {
	userRole := app.Group("/api/user-roles")
	userRole.Post("/add", controllers.CreateUserRole)
}
