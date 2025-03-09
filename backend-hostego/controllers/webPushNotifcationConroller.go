package controllers

import (
	"encoding/json"
	"fmt"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/gofiber/fiber/v3"
)

// =============================
// ✅ Push Notification Endpoint
// =============================

// VAPID keys for Web Push Notifications (Generate your own keys)
const publicKey = "BGQRMk6dwGjrQHY47G4g1gphFGBdK11REbNsz8qUkMq9XJVkLO9VWs3a72ntetjKO5PRFEyRYrWggs8VJefqr7A"
const privateKey = "W8PauXVtgDPZ8RHYulzVXEFd8uEawUwlPx8xGzMXg4w"

func SendWebPushNotification(c fiber.Ctx) error {

	var payload struct {
		Title   string `json:"title"`
		Body    string `json:"body"`
		SubJSON string `json:"subscription"`
	}

	if err := c.Bind().JSON(&payload); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Convert subscription JSON to WebPush subscription
	subscription := &webpush.Subscription{}
	if err := json.Unmarshal([]byte(payload.SubJSON), subscription); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid subscription data"})
	}

	// Create push message
	resp, err := webpush.SendNotification([]byte(fmt.Sprintf(`{"title":"%s", "body":"%s"}`, payload.Title, payload.Body)), subscription, &webpush.Options{
		VAPIDPublicKey:  publicKey,
		VAPIDPrivateKey: privateKey,
		TTL:             30,
	})

	if err != nil {
		fmt.Println("Failed to send notification:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send notification"})
	}
	defer resp.Body.Close()

	return c.JSON(fiber.Map{"message": "Notification sent successfully"})

}
