package main

import (
	"log"

	"backend-hostego/database"
	"backend-hostego/routes"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"gorm.io/gorm"
)

// Global DB variable
var db *gorm.DB

func main() {

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://hostego.in"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"Authorization"},
	}))

	database.ConnectDataBase()

	routes.AuthRoutes(app)
	routes.ShopRoutes(app)
	routes.ProductRoutes(app)
	routes.OrderRoutes(app)
	routes.WalletRoutes(app)
	routes.PaymentRoutes(app)
	routes.DeliveryPartnerRoutes(app)
	routes.CartRoutes(app)
	routes.AddressRoutes(app)
	routes.UserRoutes(app)
	routes.UserRolesRoutes(app)


	
	// Fetch all users
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the Hostego Backend Sever !"})
	})

	log.Fatal(app.Listen("0.0.0.0:8000"))
}
