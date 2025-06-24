package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api", middlewares.VerifyUserAuthCookieMiddleware())

	api.Get("/users",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "admin"),
		controllers.GetUsers,
	)

	// api.Post("/users", controllers.CreateUser)
}
