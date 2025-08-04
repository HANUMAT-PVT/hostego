package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func WebPushNotificationRoutes(app *fiber.App) {
	webPushNotification := app.Group("/api/notifications")
	webPushNotification.Post("/fcm", controllers.SendFCMNotification)
	webPushNotification.Post("/", middlewares.VerifyUserAuthCookieMiddleware(), middlewares.RoleMiddleware("super_admin"), controllers.CreateNotification)
	// webPushNotification.Post("/", controllers.SendWebPushNotification)
	webPushNotification.Get("/", middlewares.VerifyUserAuthCookieMiddleware(), middlewares.RoleMiddleware("super_admin"), controllers.GetNotifications)
}
