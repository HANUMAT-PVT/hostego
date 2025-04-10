package websocket

import (
	"log"
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
	for {
		msg := <-broadcast
		for client := range clients {
		
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
