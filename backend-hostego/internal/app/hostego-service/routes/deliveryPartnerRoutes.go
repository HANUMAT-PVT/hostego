package routes

import (
	"backend-hostego/internal/app/hostego-service/controllers"

	"github.com/gofiber/fiber/v3"
)

func DeliveryPartnerRoutes(app *fiber.App) {
	deliveryPartnerRoutes := app.Group("/api/delivery-partner")
	deliveryPartnerRoutes.Post("/", controllers.CreateNewDeliveryPartner)
	deliveryPartnerRoutes.Get("/:id", controllers.FetchDeliveryPartnerById)
}
