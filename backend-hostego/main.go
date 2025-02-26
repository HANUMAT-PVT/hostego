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
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET,POST,HEAD,PUT,DELETE,PATCH"},
		AllowHeaders:     []string{"Origin, Content-Type, Accept, Authorization, X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	database.ConnectDataBase()

	routes.SetupRoutes(app)
	routes.AuthRoutes(app)
	routes.ShopRoutes(app)
	routes.ProductRoutes(app)
	routes.OrderRoutes(app)
	routes.WalletRoutes(app)
	routes.PaymentRoutes(app)
	routes.DeliveryPartnerRoutes(app)
	routes.CartRoutes(app)
	routes.AddressRoutes(app)
	// Fetch all users
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the server Backend"})
	})

	log.Fatal(app.Listen("0.0.0.0:8080"))
}
