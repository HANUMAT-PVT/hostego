package natsclient

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gofiber/fiber/v3"
	"github.com/nats-io/nats.go"
)

// Holds active SSE connections per user
var clients = struct {
	sync.RWMutex
	connections map[string]chan string
}{
	connections: make(map[string]chan string),
}

// Polling handler for receiving messages
func PollingHandler(c fiber.Ctx) error {
	userID := c.Query("user")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user ID is required",
		})
	}

	// Ensure the message channel is created for the user if not already present
	clients.Lock()
	msgChan, exists := clients.connections[userID];
	if !exists {
		// Create a new message channel for this user
		msgChan = make(chan string, 1) // Buffer size of 1 to hold at least one message
		clients.connections[userID] = msgChan
		fmt.Println("Created new message channel for user:", userID)
	}
	clients.Unlock()

	// Wait for a message from the channel
	select {
	case msg := <-msgChan:
		// Message received for the user, send it to the frontend
		fmt.Println("Message received for user:", userID, "Message:", msg)
		return c.JSON(fiber.Map{
			"message": msg,
		})
	default:
		// No message at this moment
		fmt.Println("No message at current")
		return c.JSON(nil)
	}
}

// Start subscriber for NATS events
func StartNATSSubscriber() {
	_, err := NatsConn.Subscribe("orders.events", func(m *nats.Msg) {
		var payload map[string]string
		if err := json.Unmarshal(m.Data, &payload); err != nil {
			log.Println("Invalid NATS message:", err)
			return
		}

		userID := payload["userId"]
		if userID == "" {
			return
		}

		clients.RLock()
		if ch, ok := clients.connections[userID]; ok {
			ch <- string(m.Data)
		}
		clients.RUnlock()
	})
	if err != nil {
		log.Fatal("NATS subscription error:", err)
	}
	log.Println("✅ Subscribed to NATS events")
}

// Send message to a user
func SendMessageToUser(userID string, message string) {
	clients.RLock()
	msgChan, exists := clients.connections[userID]
	clients.RUnlock()

	if exists {
		fmt.Println("Sending message to user:", userID, "Message:", message) // Log when message is sent
		select {
		case msgChan <- message:
			fmt.Println("Message sent to channel for user:", userID)
		default:
			fmt.Println("Message channel full or closed")
		}
	} else {
		fmt.Println("No connection for user:", userID)
	}
}
