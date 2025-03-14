package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func WebPushNotificationRoutes(app *fiber.App) {
	webPushNotification := app.Group("/api/notifications")
	webPushNotification.Post("/", controllers.SendWebPushNotification)
}
