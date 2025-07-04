package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v2"
)

func WebPushNotificationRoutes(app *fiber.App) {
	webPushNotification := app.Group("/api/notifications")
	webPushNotification.Post("/fcm", controllers.SendFCMNotification)
	// webPushNotification.Post("/", controllers.SendWebPushNotification)
}
