package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v3"
)

func UserRolesRoutes(app *fiber.App) {
	userRoles := app.Group("/api/user-roles", middlewares.VerifyUserAuthCookieMiddleware())
	userRoles.Get("/", controllers.FetchUserRoles)

	userRoles.Post("/add", middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "admin"), controllers.CreateUserRole,
	)

	userRoles.Delete("/:id",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "admin"),
		controllers.DeleteUserRole,
	)
}
