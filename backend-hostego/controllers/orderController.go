package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"
	"encoding/json"

	"math"

	"time"

	"github.com/gofiber/fiber/v3"
)

type FinalOrderValueType struct {
	SubTotal             float64 `json:"subtotal"`
	ShippingFee          float64 `json:"shipping_fee"`
	PlatformFee          float64 `json:"platform_fee"`
	DeliveryPartnerShare float64 `json:"delivery_partner_fee"`
	FinalOrderValue      float64 `json:"final_order_value"`
}

// Move struct definition outside the function
type requestCreateOrder struct {
	AddressId string `json:"address_id"`
}

func CreateNewOrder(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	var cartItems []models.CartItem
	var order models.Order
	var requestOrder requestCreateOrder // Use the struct here
	if err := c.Bind().JSON(&requestOrder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if requestOrder.AddressId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Address ID is required"})
	}

	if err := database.DB.Preload("ProductItem.Shop").Where("user_id=?", user_id).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	jsonCartItems, err := json.Marshal(cartItems)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to serialize cart items"})
	}
	order.OrderItems = jsonCartItems
	order.UserId = user_id
	totalCharges := CalculateFinalOrderValue(cartItems)
	order.PlatformFee = totalCharges.PlatformFee
	order.ShippingFee = totalCharges.ShippingFee
	order.FinalOrderValue = totalCharges.FinalOrderValue
	order.DeliveryPartnerFee = totalCharges.DeliveryPartnerShare
	order.OrderStatus = "pending"
	order.AddressID = requestOrder.AddressId
	if err := database.DB.Preload("User").Create(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// CalculateFinalOrderValue calculates the total order cost including charges
func CalculateFinalOrderValue(cartItems []models.CartItem) FinalOrderValueType {
	totalItemSubTotal := 0.0

	// Calculate the subtotal of all cart items
	for _, item := range cartItems {
		totalItemSubTotal = item.SubTotal + totalItemSubTotal
	}

	// Calculate charges
	var deliveryCharge float64
	if totalItemSubTotal <= 150.0 {
		deliveryCharge = 15.0
	} else {
		charge := totalItemSubTotal * 0.10
		if charge > 30.0 {
			deliveryCharge = 30.0
		} else {
			deliveryCharge = math.Round(charge*100) / 100 // Round to 2 decimal places
		}
	}

	// Platform fee (fixed at ₹1)
	platformFee := 1.0

	// Distribution of charge (70% for delivery partner, 30% for company)
	deliveryPartnerShare := math.Round((deliveryCharge*0.70)*100) / 100

	// Final order value including charges
	finalOrderValue := math.Round((totalItemSubTotal+deliveryCharge)*100) / 100

	return FinalOrderValueType{
		SubTotal:             math.Round(totalItemSubTotal*100) / 100,
		ShippingFee:          deliveryCharge,
		PlatformFee:          platformFee,
		DeliveryPartnerShare: deliveryPartnerShare,
		FinalOrderValue:      finalOrderValue,
	}
}

func MarkOrderAsDelivered(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if user_id == "" {

	}
	orderId := c.Params("id")
	var order models.Order
	if err := database.DB.First(&order, "order_id = ?", orderId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Set delivered time to current time
	order.DeliveredAt = time.Now()

	// Save changes
	if err := database.DB.Save(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"order": order, "message": "Order is delivered succesfully!"})
}

func FetchOrderById(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if user_id == "" {

	}
	order_id := c.Params("id")
	var order models.Order

	if err := database.DB.Preload("User").Preload("PaymentTransaction").Preload("Address").Where("order_id=?", order_id).First(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

func AssignOrderToDeliveryPartner(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)

	type requestAssignOrder struct {
		DeliveryPartnerId string `json:"delivery_partner_id"`
		OrderId           string `json:"order_id"`
	}
	var request_assign requestAssignOrder

	var order models.Order

	var delivery_partner models.DeliveryPartner

	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if user_id == "" {

	}
	if err := c.Bind().JSON(&request_assign); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := database.DB.Where("order_id=?", request_assign.OrderId).Find(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if err := database.DB.Preload("User").Where("delivery_partner_id=?", request_assign.DeliveryPartnerId).Find(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	jsonDeliveryPartner, err := json.Marshal(delivery_partner)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	order.DeliveryPartner = jsonDeliveryPartner
	order.DeliveryPartnerId = request_assign.DeliveryPartnerId
	order.OrderStatus = models.AssignedOrderStatus

	if err := database.DB.Where("order_id=?", request_assign.OrderId).Save(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order assigned succefully to" + " " + delivery_partner.User.FirstName})

}

func UpdateOrderById(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if user_id == "" {
	}
	order_id := c.Params("id")

	var existingOrder models.Order
	if err := database.DB.Where("order_id = ?", order_id).First(&existingOrder).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	var updateData struct {
		OrderStatus models.OrderStatusType `json:"order_status"`
	}
	if err := c.Bind().JSON(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Update only the status
	existingOrder.OrderStatus = updateData.OrderStatus

	// Save the changes
	if err := database.DB.Save(&existingOrder).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Order status updated successfully!",
		"order":   existingOrder,
	})
}

func FetchAllUserOrders(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if user_id == "" {

	}
	var orders []models.Order
	if err := database.DB.Preload("User").Preload("PaymentTransaction").Preload("Address").Where("user_id=?", user_id).Order("created_at desc").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(orders)
}

func FetchAllOrders(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if user_id == "" {
	}

	dbQuery := database.DB

	searchQuery := c.Query("search")

	filter := c.Query("filter")

	if searchQuery != "" {
		dbQuery = dbQuery.Where(
			"order_id LIKE ? OR user_id IN (SELECT user_id FROM users WHERE mobile_number LIKE ?)",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
		)
	}

	if filter != "" {
		dbQuery = dbQuery.Where("order_status = ?", filter)
	}

	var orders []models.Order
	if err := dbQuery.Preload("User").Preload("PaymentTransaction").Preload("Address").Order("created_at desc").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

func FetchAllOrdersByDeliveryPartner(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}

	if user_id == "" {
	}
	// Get delivery partner ID from params
	delivery_partner_id := c.Params("id")

	// Initialize query with preloads
	dbQuery := database.DB.Preload("User").
		Preload("PaymentTransaction").
		Preload("Address")

	// Get query parameters for filtering
	status := c.Query("status")
	searchQuery := c.Query("search")

	// Base query for delivery partner's orders
	dbQuery = dbQuery.Where("delivery_partner_id = ?", delivery_partner_id)

	// Add status filter if provided
	if status != "" {
		dbQuery = dbQuery.Where("order_status = ?", status)
	}

	// Add search functionality
	if searchQuery != "" {
		dbQuery = dbQuery.Where(
			"order_id LIKE ? OR user_id IN (SELECT user_id FROM users WHERE mobile_number LIKE ?)",
			"%"+searchQuery+"%",
			"%"+searchQuery+"%",
		)
	}

	var orders []models.Order
	if err := dbQuery.Order("created_at desc").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"orders": orders,
		"count":  len(orders),
	})
}
