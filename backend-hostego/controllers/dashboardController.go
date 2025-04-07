package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func GetDashBoardStats(c fiber.Ctx) error {
	var result struct {
		TotalCustomers         int     `json:"total_customers"`
		RepeatCustomers        int     `json:"repeat_customers"`
		RepeatCustomersRevenue float64 `json:"repeat_customers_revenue"`
	}

	var totalShipping float64
	var totalDeliveryFee float64

	// Total Shipping Fee
	err := database.DB.Model(&models.Order{}).
		Select("SUM(shipping_fee)").
		Where("order_status = ?", "delivered").
		Scan(&totalShipping).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch shipping fee"})
	}

	// Total Delivery Partner Fee
	err = database.DB.Model(&models.Order{}).
		Select("SUM(delivery_partner_fee)").
		Where("order_status = ?", "delivered").
		Scan(&totalDeliveryFee).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch delivery partner fee"})
	}

	// Repeat customer stats via raw SQL
	query := `
	SELECT
	  COUNT(DISTINCT o.user_id) AS total_customers,
	  COUNT(DISTINCT CASE WHEN user_order_count > 1 THEN o.user_id END) AS repeat_customers,
	  SUM(CASE WHEN user_order_count > 1 THEN o.final_order_value ELSE 0 END) AS repeat_customers_revenue
	FROM (
	  SELECT user_id, COUNT(*) AS user_order_count
	  FROM orders
	  GROUP BY user_id
	) AS user_orders
	JOIN orders o ON o.user_id = user_orders.user_id;
	`
	err = database.DB.Raw(query).Scan(&result).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch customer stats"})
	}

	return c.JSON(fiber.Map{
		"result":                     result,
		"total_shipping_revenue":     totalShipping,
		"total_delivery_partner_fee": totalDeliveryFee,
	})
}
