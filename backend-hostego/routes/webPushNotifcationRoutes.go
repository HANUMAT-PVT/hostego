package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func WebPushNotificationRoutes(app *fiber.App) {
	webPushNotification := app.Group("/api/notifications")
	webPushNotification.Post("/fcm", middlewares.VerifyUserAuthCookieMiddleware(), middlewares.RoleMiddleware("admin", "super_admin"), controllers.SendFCMNotification)
	// webPushNotification.Post("/", controllers.SendWebPushNotification)
}
