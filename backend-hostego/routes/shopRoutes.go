package routes

import (
	"backend-hostego/controllers"
	"backend-hostego/middlewares"
	"github.com/gofiber/fiber/v3"
)

func ShopRoutes(app *fiber.App) error {

	shopRoutes:=	app.Group("/api/shop", middlewares.VerifyUserAuthCookieMiddleware())

	shopRoutes.Post("/", controllers.CreateShop)

	shopRoutes.Get("/:id", controllers.FetchShopById)

	shopRoutes.Get("/",controllers.FetchShops)

	shopRoutes.Patch("/:id", controllers.UpdateShopById)

	return nil
}

// zindagi ki raah pr hum yu  khade hai mayusiyat ko peeche chod hum aage badhe