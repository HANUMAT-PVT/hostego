package routes

import (
	"backend-hostego/controllers"

	"github.com/gofiber/fiber/v3"
)

func DeliveryPartnerRoutes(app *fiber.App) {
	deliveryPartnerRoutes := app.Group("/api/delivery-partner")
	deliveryPartnerRoutes.Post("/", controllers.CreateNewDeliveryPartner)
	deliveryPartnerRoutes.Get("/find", controllers.FetchDeliveryPartnerByUserId)
	deliveryPartnerRoutes.Patch("/:id", controllers.UpdateDeliveryPartner)
}
