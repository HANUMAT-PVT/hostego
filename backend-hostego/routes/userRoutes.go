package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func UserRoutes(app *fiber.App) {
	userRoutes := app.Group("/api/users", middlewares.VerifyUserAuthCookieMiddleware())
	userRoutes.Patch("/me", controllers.UpdateUserById)
	userRoutes.Get("/me", controllers.GetUserById)
	userRoutes.Get("/", controllers.GetUsers)
}
