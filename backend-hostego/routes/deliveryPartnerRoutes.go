package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v2"
)

func DeliveryPartnerRoutes(app *fiber.App) {
	deliveryPartnerRoutes := app.Group("/api/delivery-partner", middlewares.VerifyUserAuthCookieMiddleware())

	deliveryPartnerRoutes.Get("/all",
		middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "delivery_partner_manager"),
		controllers.FetchAllDeliveryPartners,
	)

	deliveryPartnerRoutes.Post("/",controllers.CreateNewDeliveryPartner)

	deliveryPartnerRoutes.Get("/find", controllers.FetchDeliveryPartnerByUserId)
	deliveryPartnerRoutes.Patch("/:id", controllers.UpdateDeliveryPartner)
	deliveryPartnerRoutes.Get("/earnings/:id", controllers.FetchDeliveryPartnerEarnings)
}
