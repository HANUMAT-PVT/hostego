package services

import (
	"backend-hostego/models"
	"context"
	"errors"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/api/option"
)

var Client *messaging.Client

// Load once at startup
func Init(ctx context.Context, credentialFile string) error {
	var pathForFirebase = "/etc/hostego/firebase.json"
	//var pathForFirebase = credentialFile

	opt := option.WithCredentialsFile(pathForFirebase)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}
	Client, err = app.Messaging(ctx)
	if err == nil {
		log.Println("✅  FCM initialized")
	}
	return err
}

// One‑shot send helper that returns message ID and error
func SendToToken(ctx context.Context, token, title, body string, data map[string]string) (string, error) {
	if Client == nil {
		return "", errors.New("fcm client not initialized")
	}

	msg := &messaging.Message{
		Token: token,
		Data:  data,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Android: &messaging.AndroidConfig{Priority: "high"},
		APNS: &messaging.APNSConfig{Headers: map[string]string{
			"apns-priority": "10",
		}},
	}

	// Send the message and capture the message ID
	msgID, err := Client.Send(ctx, msg)
	if err != nil {
		return "", err
	}

	return msgID, nil
}

// Lift's Down? Let Us Lift Your Mood!
// Don’t climb hungry — we’ll bring your food right to your floor.
// // 🍔 Order now & chill!
// ⏰ Just 30 minutes left!
// Your favorite kitchen is about to close — order now before it’s too late!

func SendFCMNotification(c *fiber.Ctx) error {

	var notification models.Notification

	if err := c.BodyParser(&notification); err != nil {
		return err
	}

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title:    notification.Title,
			Body:     notification.Body,
			ImageURL: notification.NotificationImageUrl,
		},
		Topic: "hostego_updates",
	}

	_, err := Client.Send(context.Background(), message)
	if err != nil {
		log.Fatalf("❌ error sending message: %v", err)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Notification sent successfully !"})
}
