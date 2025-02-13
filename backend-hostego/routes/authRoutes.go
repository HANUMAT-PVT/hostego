package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func AuthRoutes(app *fiber.App){
	auth:=app.Group("/auth")
	auth.Post("/signup",controllers.Signup);
}