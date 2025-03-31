package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"

	"github.com/gofiber/fiber/v3"
)

func ShopRoutes(app *fiber.App) {

	shopRoutes := app.Group("/api/shop", middlewares.VerifyUserAuthCookieMiddleware(),
		middlewares.RoleMiddleware("super_admin", "admin", "inventory_manager"))

	shopRoutes.Post("/", controllers.CreateShop)

	shopRoutes.Get("/:id", controllers.FetchShopById)

	shopRoutes.Get("/", controllers.FetchShops)

	shopRoutes.Patch("/:id", controllers.UpdateShopById)

}

// zindagi ki raah pr hum yu  khade hai mayusiyat ko peeche chod hum aage badhe
