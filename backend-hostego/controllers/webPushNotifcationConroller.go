package controllers

import (
	"encoding/json"
	"fmt"
	"os"

	webpush "github.com/SherClockHolmes/webpush-go"
	"github.com/gofiber/fiber/v2"
)

// =============================
// ✅ Push Notification Endpoint
// =============================

// VAPID keys for Web Push Notifications (Generate your own keys)
const publicKey = "BGQRMk6dwGjrQHY47G4g1gphFGBdK11REbNsz8qUkMq9XJVkLO9VWs3a72ntetjKO5PRFEyRYrWggs8VJefqr7A"
const privateKey = "W8PauXVtgDPZ8RHYulzVXEFd8uEawUwlPx8xGzMXg4w"

func SendWebPushNotification(c *fiber.Ctx) error {

	var payload struct {
		Title   string `json:"title"`
		Body    string `json:"body"`
		SubJSON string `json:"subscription"`
	}

	if err := c.BodyParser(&payload); err != nil {
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

// SendWebPushNotification sends a push notification with the given title and body
func CreateWebPushNotification(title, body string) error {
	// Step 1: Read subscription data from a file
	data, err := os.ReadFile("subscriptions.json")

	if err != nil {
		return fmt.Errorf("failed to read subscription file: %w", err)
	}
	println("enterin gthis")

	// Step 2: Parse JSON into webpush.Subscription struct
	subscription := &webpush.Subscription{}
	if err := json.Unmarshal(data, subscription); err != nil {
		return fmt.Errorf("failed to unmarshal subscription JSON: %w", err)
	}

	// Step 3: Create the message payload
	message := map[string]string{
		"title": "hHello",
		"body":  "launching",
	}
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}

	// Step 4: Send the push notification
	resp, err := webpush.SendNotification(messageJSON, subscription, &webpush.Options{
		Subscriber:      "mailto:admin@example.com", // optional but recommended
		VAPIDPublicKey:  publicKey,
		VAPIDPrivateKey: privateKey,
		TTL:             30,
	})
	if err != nil {
		return fmt.Errorf("failed to send push notification: %w", err)
	}
	defer resp.Body.Close()

	// Optional: Check response status
	fmt.Println("Notification sent with status:", resp.Status)

	return nil
}
