package natsclient

import (
	"log"

	"github.com/nats-io/nats.go"
)

var NatsConn *nats.Conn

func ConnectNATS() {
	var err error
	NatsConn, err = nats.Connect("nats://localhost:4222") // Use your actual NATS server URL
	if err != nil {
		log.Fatal("Failed to connect to NATS:", err)
	}
	log.Println("✅ Connected to NATS")
}
