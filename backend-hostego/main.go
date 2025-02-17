package main

import (
	"log"

	"backend-hostego/database"
	"backend-hostego/routes"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// Global DB variable
var db *gorm.DB

func main() {

	app := fiber.New()

	database.ConnectDataBase()

	routes.SetupRoutes(app)
	routes.AuthRoutes(app)
	routes.ShopRoutes(app)
	routes.ProductRoutes(app)
	routes.OrderRoutes(app)
	routes.CartRoutes(app)
	routes.WalletRoutes(app)
	routes.PaymentRoutes(app)
	routes.DeliveryPartnerRoutes(app)
	// Fetch all users
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the server Backend"})
	})

	log.Fatal(app.Listen(":8080"))
}
