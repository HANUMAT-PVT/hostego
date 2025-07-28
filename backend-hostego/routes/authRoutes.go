package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	auth := app.Group("/api/auth")
	auth.Post("/signup", controllers.Signup)
	auth.Delete("/delete-user", middlewares.VerifyUserAuthCookieMiddleware(), controllers.DeleteUserById)
}
