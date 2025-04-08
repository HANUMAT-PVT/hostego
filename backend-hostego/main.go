package main

import (
	"backend-hostego/cron"
	"backend-hostego/database"
	"backend-hostego/routes"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	websocket "github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var db *gorm.DB

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
	}
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			log.Println(err.Error())
			break
		}
		// print the recived message from the frontned
		log.Println("Recieved message %s", message)
		// send message to the clients
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println(err.Error())
			break
		}
	}
}

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

	http.HandleFunc("/websocket", websocketHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))

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
	// Default Route
	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to the Hostego Backend Server!"})
	})

	// Initialize cron jobs
	cron.InitCronJobs()

	// Start the server
	log.Fatal(app.Listen("0.0.0.0:8000"))
}
