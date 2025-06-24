package natsclient

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
)

// Holds active SSE connections per user
type ClientConnection struct {
	MsgChan chan string
	Roles   []string
}

var clients = struct {
	sync.RWMutex
	connections map[string]ClientConnection
}{
	connections: make(map[string]ClientConnection),
}

// Polling handler for receiving messages
func PollingHandler(c *fiber.Ctx) error {
	userID := c.Query("user")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Get roles from query param: ?roles=admin,order_assign_manager
	rolesParam := c.Query("roles")
	roles := strings.Split(rolesParam, ",")

	// Sanitize empty string role
	if len(roles) == 1 && roles[0] == "" {
		roles = []string{}
	}

	clients.Lock()
	conn, exists := clients.connections[userID]
	if !exists {
		// New connection
		conn = ClientConnection{
			MsgChan: make(chan string, 1),
			Roles:   roles,
		}
		fmt.Println("Created new connection for user:", userID, "Roles:", roles)
	} else {
		// Existing connection, update roles if needed
		conn.Roles = roles
	}
	clients.connections[userID] = conn
	clients.Unlock()

	// Wait for a message from the channel
	select {
	case msg := <-conn.MsgChan:
		fmt.Println("Message received for user:", userID, "Message:", msg)
		return c.JSON(fiber.Map{
			"message": msg,
		})
	default:
		fmt.Println("No message currently for user:", userID)
		return c.JSON(nil)
	}
}

// Start subscriber for NATS events
func StartNATSSubscriber() {
	_, err := NatsConn.Subscribe("orders.events", func(m *nats.Msg) {
		var payload map[string]interface{}
		if err := json.Unmarshal(m.Data, &payload); err != nil {
			log.Println("Invalid NATS message:", err)
			return
		}

		// Get the roles from the payload
		roles, exists := payload["roles"].([]interface{})
		if !exists || len(roles) == 0 {
			return
		}

		// Lock the clients map to safely access it
		clients.RLock()
		// Iterate through the connections and send messages to users with matching roles
		for userID, conn := range clients.connections {
			// Check if the user has at least one of the roles
			for _, role := range conn.Roles {
				for _, r := range roles {
					if role == r {
						// Send the message to the user's channel if their role matches
						select {
						case conn.MsgChan <- string(m.Data): // Send the message to the channel
							fmt.Println("Message sent to user:", userID, "Roles:", conn.Roles)
						default:
							fmt.Println("Message channel full or closed for user:", userID)
						}
					}
				}
			}
		}
		clients.RUnlock()
	})

	if err != nil {
		log.Fatal("NATS subscription error:", err)
	}
	log.Println("✅ Subscribed to NATS events")
}

type UserMessage struct {
	Title string   `json:"title"`
	Body  string   `json:"body"`
	Roles []string `json:"roles"` // Slice to hold multiple roles
}

func SendMessageToUsersByRole(roles []string, title string, body string) {
	clients.RLock()
	// Iterate over all connections and find users with matching roles
	for userID, conn := range clients.connections {
		// Check if the user has at least one matching role
		for _, role := range conn.Roles {
			for _, r := range roles {
				if role == r {
					// Send the message to this user's channel
					userMessage := UserMessage{
						Title: title,
						Body:  body,
						Roles: conn.Roles,
					}

					jsonMessage, err := json.Marshal(userMessage)
					if err != nil {
						fmt.Println("Error marshalling message:", err)
						return
					}

					select {
					case conn.MsgChan <- string(jsonMessage): // ✅ Now we are sending to the channel inside the struct
						fmt.Println("Message sent to channel for user:", userID)
					default:
						fmt.Println("Message channel full or closed for user:", userID)
					}
				}
			}
		}
	}
	clients.RUnlock()
}
