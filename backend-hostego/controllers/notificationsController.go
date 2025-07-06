package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"backend-hostego/services"
	"context"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"
)

func NotifyOrderPlaced(orderID int) error {
	db := database.DB
	// 1. Get order
	var order models.Order
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

	// 4. Get Shop Owner using OwnerId
	var owner models.User
	if err := db.First(&owner, shop.OwnerId).Error; err != nil {
		return err
	}

	// 5. Send notifications
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var g errgroup.Group

	// Notification to Customer
	if customer.FCMToken != "" {
		g.Go(func() error {
			_, err := services.SendToToken(
				context.Background(),
				customer.FCMToken,
				"Payment Successful",
				"Your payment for Order #"+strconv.Itoa(order.OrderId)+" was successful 🎉",
				map[string]string{
					"type":     "new_order",
					"order_id": strconv.Itoa(order.OrderId),
				},
			)
			return err
		})
	}

	// Notification to Restaurant Owner
	if owner.FCMToken != "" {
		g.Go(func() error {
			_, err := services.SendToToken(
				ctx,
				owner.FCMToken,
				"New Order Received",
				"Order #"+strconv.Itoa(order.OrderId)+" has been placed. Please confirm it.",
				map[string]string{
					"type":     "new_order_request",
					"order_id": strconv.Itoa(order.OrderId),
				},
			)
			return err
		})
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

	var payload map[string]any

	payload = map[string]any{
		"title": map[bool]string{true: "Order Accepted", false: "Order Rejected"}[isAccepted],
		"body":  map[bool]string{true: "Your order #" + strconv.Itoa(order.OrderId) + " has been accepted by the restaurant. Expected ready in " + strconv.Itoa(expectedReadyInMins) + " mins.", false: "Your order #" + strconv.Itoa(order.OrderId) + " has been rejected by the restaurant."}[isAccepted],
		"content": map[string]string{
			"type":     "new_order",
			"order_id": strconv.Itoa(order.OrderId),
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
				payload["content"].(map[string]string),
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

	if customer.FCMToken != "" {
		services.SendToToken(context.Background(), customer.FCMToken, title, message, map[string]string{"order_id": strconv.Itoa(order.OrderId), "type": "new_order"})
	}

	return nil
}
