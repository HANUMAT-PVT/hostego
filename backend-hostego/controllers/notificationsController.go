package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"backend-hostego/services"
	"context"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

func NotifyOrderPlaced(orderID int) error {
	db := database.DB
	// 1. Get order
	var order models.Order
	var owner models.User
	var userRoles []models.UserRole

	if err := db.Table("user_roles").Preload("User").Where("role_id = ?", 6).Find(&userRoles).Error; err != nil {
		return err
	}

	if err := db.First(&order, orderID).Error; err != nil {
		return err
	}

	// 2. Get Shop (to get OwnerId)
	var shop models.Shop
	if err := db.First(&shop, order.ShopId).Error; err != nil {
		return err
	}

	// 3. Get User (Customer)
	var customer models.User
	if err := db.First(&customer, order.UserId).Error; err != nil {
		return err
	}

	// 5. Send notifications
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var g errgroup.Group

	payload := map[string]any{
		"title": "Payment Successful",
		"body":  "Your payment for Order #" + strconv.Itoa(order.OrderId) + " was successful 🎉",
		"data": map[string]string{
			"type":               "new_order",
			"order_id":           strconv.Itoa(order.OrderId),
			"android_channel_id": "order_updates",
		},
	}
	// Notification to Customer
	if customer.FCMToken != "" {
		g.Go(func() error {
			_, err := services.SendToToken(
				context.Background(),
				customer.FCMToken,
				payload["title"].(string),
				payload["body"].(string),
				payload["data"].(map[string]string),
			)
			return err
		})
	}

	// Notification to Restaurant Owner
	fmt.Println("shop.OwnerId", shop.OwnerId)

	if shop.OwnerId != 0 {
		// 4. Get Shop Owner using OwnerId

		if err := db.First(&owner, shop.OwnerId).Error; err != nil {
			return err
		}

		if owner.FCMToken != "" {
			g.Go(func() error {
				_, err := services.SendToToken(
					ctx,
					owner.FCMToken,
					"New Order Received",
					"Order #"+strconv.Itoa(order.OrderId)+" has been placed. Please confirm it.",
					map[string]string{
						"type":               "new_order_request",
						"order_id":           strconv.Itoa(order.OrderId),
						"android_channel_id": "order_updates",
					},
				)
				return err
			})
		}
	}

	payload = map[string]any{
		"type":  "new_order_request",
		"title": "New Order Request",
		"body":  "Order #" + strconv.Itoa(order.OrderId) + " has been placed. Please assign it to a delivery partner.",
		"data": map[string]string{
			"type":               "new_order_request",
			"order_id":           strconv.Itoa(order.OrderId),
			"android_channel_id": "order_updates",
		},
	}

	for _, userRole := range userRoles {
		if userRole.User.FCMToken != "" {
			g.Go(func() error {
				_, err := services.SendToToken(ctx, userRole.User.FCMToken,
					payload["title"].(string),
					payload["body"].(string),
					payload["data"].(map[string]string),
				)
				return err
			})
		}
	}

	return g.Wait()
}

func NotifyOrderAcceptedOrRejectedByRestaurant(orderID int, isAccepted bool, expectedReadyInMins int) error {
	db := database.DB
	// 1. Get order
	var order models.Order
	if err := db.First(&order, orderID).Error; err != nil {
		return err
	}

	// 3. Get User (Customer)
	var customer models.User
	if err := db.First(&customer, order.UserId).Error; err != nil {
		return err
	}

	// 5. Send notifications
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var g errgroup.Group

	payload := map[string]any{
		"title": map[bool]string{true: "Order Accepted", false: "Order Rejected"}[isAccepted],
		"body":  map[bool]string{true: "Your order #" + strconv.Itoa(order.OrderId) + " has been accepted by the restaurant. Expected ready in " + strconv.Itoa(expectedReadyInMins) + " mins.", false: "Your order #" + strconv.Itoa(order.OrderId) + " has been rejected by the restaurant."}[isAccepted],
		"data": map[string]string{
			"type":               "new_order",
			"order_id":           strconv.Itoa(order.OrderId),
			"android_channel_id": "order_updates",
		},
	}

	// Notification to Customer
	if customer.FCMToken != "" {
		g.Go(func() error {
			_, err := services.SendToToken(
				ctx,
				customer.FCMToken,
				payload["title"].(string),
				payload["body"].(string),
				payload["data"].(map[string]string),
			)
			return err
		})
	}

	return g.Wait()

}

// func NotifyOrderDelivered(orderID int) error {
// 	db := database.DB
// 	// 1. Get order
// 	var order models.Order
// 	if err := db.First(&order, orderID).Error; err != nil {
// 		return err
// 	}
// }

func NotifyOrderToCustomerByRestaurant(orderID int, message string, title string) error {
	db := database.DB
	// 1. Get order
	var order models.Order
	if err := db.First(&order, orderID).Error; err != nil {
		return err
	}

	var customer models.User
	if err := db.First(&customer, order.UserId).Error; err != nil {
		return err
	}

	payload := map[string]any{
		"title": title,
		"body":  message,
		"data": map[string]string{
			"type":               "new_order",
			"order_id":           strconv.Itoa(order.OrderId),
			"android_channel_id": "order_updates",
			"payload":            "new_order:" + strconv.Itoa(order.OrderId),
		},
	}

	if customer.FCMToken != "" {
		_, err := services.SendToToken(context.Background(),
			customer.FCMToken,
			payload["title"].(string),
			payload["body"].(string),
			payload["data"].(map[string]string),
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// this is used to notify a person by their user id and order id
// whether they are customer or restaurant owner or delivery partner
func NotifyPersonByUserIdAndOrderID(orderID int, message string, title string, userID int) error {
	db := database.DB
	// 1. Get order
	var order models.Order
	if err := db.First(&order, orderID).Error; err != nil {
		return err
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		return err
	}

	payload := map[string]any{
		"title": title,
		"body":  message,
		"data": map[string]string{
			"type":               "new_order",
			"order_id":           strconv.Itoa(order.OrderId),
			"android_channel_id": "order_updates",
		},
	}

	if user.FCMToken != "" {
		_, err := services.SendToToken(
			context.Background(),
			user.FCMToken,
			payload["title"].(string),
			payload["body"].(string),
			payload["data"].(map[string]string),
		)
		if err != nil {
			return err
		}
	}

	return nil
}
