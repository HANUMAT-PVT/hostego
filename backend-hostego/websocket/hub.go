package websocket

import (
	"log"
	"runtime/debug"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Role string // You can replace this with userID or role
}

var clients = make(map[*Client]bool)
var broadcast = make(chan Message)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func HandleMessages() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("🚨 CRITICAL: WebSocket handler panic: %v", r)
			log.Printf("Stack trace: %s", debug.Stack())
		}
	}()

	for {
		msg := <-broadcast
		// Create a safe copy of clients to avoid concurrent map iteration issues
		clientsCopy := make(map[*Client]bool)
		for client, active := range clients {
			clientsCopy[client] = active
		}

		for client := range clientsCopy {
			if client.Role == msg.Role {
				err := client.Conn.WriteJSON(msg)
				if err != nil {
					log.Printf("WebSocket error: %v", err)
					client.Conn.Close()
					delete(clients, client)
				}
			}
		}
	}
}

func RegisterClient(client *Client) {
	clients[client] = true
}

func SendMessage(msg Message) {
	broadcast <- msg
}
