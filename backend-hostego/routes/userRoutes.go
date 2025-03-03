package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func UserRoutes(app *fiber.App) {
	userRoutes := app.Group("/api/user")
	userRoutes.Patch("/me", controllers.UpdateUserById)
}
