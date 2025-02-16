package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"
	"time"

	"github.com/gofiber/fiber/v3"
)

func CreateNewOrder(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	var cartItems []models.CartItem

	var order models.Order
	if err := database.DB.Preload("Product").Where("user_id=?", user_id).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if err := c.Bind().JSON(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	order.OrderItems = cartItems
	totalCharges := CalculateFinalOrderValue(cartItems)
	order.PlatformFee = totalCharges.PlatformFee
	order.ShippingFee = totalCharges.ShippingFee
	order.FinalOrderValue = totalCharges.FinalOrderValue
	order.DeliveryPartnerFee = totalCharges.DeliveryPartnerShare

	if err := database.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"order": order, "message": "Order created successfully !"})
}

type FinalOrderValueType struct {
	SubTotal             float64 `json:"subtotal"`
	ShippingFee          float64 `json:"shipping_fee"`
	PlatformFee          float64 `json:"platform_fee"`
	DeliveryPartnerShare float64 `json:"delivery_partner_fee"`
	FinalOrderValue      float64 `json:"final_order_value"`
}

// CalculateFinalOrderValue calculates the total order cost including charges
func CalculateFinalOrderValue(cartItems []models.CartItem) FinalOrderValueType {
	totalItemSubTotal := 0.0

	// Calculate the subtotal of all cart items
	for _, item := range cartItems {
		totalItemSubTotal = item.SubTotal
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
			deliveryCharge = charge
		}
	}

	// Platform fee (fixed at ₹1)
	platformFee := 1.0

	// Distribution of charge (70% for delivery partner, 30% for company)
	deliveryPartnerShare := deliveryCharge * 0.70

	// Final order value including charges
	finalOrderValue := totalItemSubTotal + deliveryCharge + platformFee

	return FinalOrderValueType{
		SubTotal:             totalItemSubTotal,
		ShippingFee:          deliveryCharge,
		PlatformFee:          platformFee,
		DeliveryPartnerShare: deliveryPartnerShare,
		FinalOrderValue:      finalOrderValue}
}

func MarkOrderAsDelivered(c fiber.Ctx) error {
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
	order_id := c.Params("id")
	var order models.Order

	if err := database.DB.Where("order_id=?", order_id).First(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"order": order})
}
