package main

import (
	"backend-hostego/cron"
	"backend-hostego/database"
	"backend-hostego/routes"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"

	"gorm.io/gorm"
)

var db *gorm.DB

// VAPID keys for Web Push Notifications (Generate your own keys)
const publicKey = "BGQRMk6dwGjrQHY47G4g1gphFGBdK11REbNsz8qUkMq9XJVkLO9VWs3a72ntetjKO5PRFEyRYrWggs8VJefqr7A"
const privateKey = "W8PauXVtgDPZ8RHYulzVXEFd8uEawUwlPx8xGzMXg4w"

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://hostego.in"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"Authorization"},
	}))

	// Connect Database
	database.ConnectDataBase()

	// Setup All Routes
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
	routes.WebPushNotificationRoutes(app)
	routes.SearchQueryRoutes(app)
	routes.DeliveryPartnerWalletRoutes(app)
	routes.MessMenuRoutes(app)
	// Default Route
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the Hostego Backend Server!"})
	})

	// Initialize cron jobs
	cron.InitCronJobs()

	// Start the server
	log.Fatal(app.Listen("0.0.0.0:8000"))
}
