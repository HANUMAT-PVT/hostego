package main

import (
	"backend-hostego/cron"
	"backend-hostego/database"
	"backend-hostego/routes"
	"backend-hostego/services"
	"context"
	"time"

	websocket "github.com/gofiber/websocket/v2"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

var db *gorm.DB

// VAPID keys for Web Push Notifications (Generate your own keys)
const publicKey = "BGQRMk6dwGjrQHY47G4g1gphFGBdK11REbNsz8qUkMq9XJVkLO9VWs3a72ntetjKO5PRFEyRYrWggs8VJefqr7A"
const privateKey = "W8PauXVtgDPZ8RHYulzVXEFd8uEawUwlPx8xGzMXg4w"

func main() {

	app := fiber.New(fiber.Config{
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, https://hostego.in,https://www.hostego.in",
		AllowMethods:     "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: false,
		ExposeHeaders:    "Authorization",
	}))

	app.Use(recover.New()) // Logs panics
	app.Use(logger.New())  // Logs requests: method, path, status, latency
	// Connect Database
	database.ConnectDataBase()

	// natsclient.ConnectNATS()
	// natsclient.StartNATSSubscriber()

	// Add SSE route

	// WebSocket route
	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		for {
			// Read message from client
			mt, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				break
			}
			log.Printf("recv: %s", msg)

			// Echo the message back
			err = c.WriteMessage(mt, msg)
			if err != nil {
				log.Println("write error:", err)
				break
			}
		}
	}))

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
	routes.OrderItemRoutes(app)
	routes.DashboardRoutes(app)
	routes.PaymentWebhookRoutes(app)
	routes.ProductCategoryRoutes(app)
	routes.RatingRoutes(app)
	routes.ShopDashboardRoutes(app)
	// Default Route

	if err := services.Init(context.Background(), "config/firebase-service-account.json"); err != nil {
		log.Fatal("FCM init failed:", err)
	} else {
		log.Println("FCM initialized")
	}
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the Hostego Backend Server!"})
	})
	// app.Get("/events", natsclient.PollingHandler)

	// Initialize cron jobs
	cron.InitCronJobs()

	// logs.InitLogger()

	// Start the server
	log.Fatal(app.Listen("0.0.0.0:8000"))
	// firebase service account

}
