package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func RoleRoutes(app *fiber.App) {
	roleRoutes := app.Group("/api/roles", middlewares.VerifyUserAuthCookieMiddleware())
	roleRoutes.Get("/", controllers.FetchUserRoles)
	roleRoutes.Delete("/:user_role_id", controllers.DeleteUserRole)
}
