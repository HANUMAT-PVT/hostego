package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetDashBoardStats(c *fiber.Ctx) error {
	// Get query parameters for date filtering
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	// Parse start date
	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid start_date format. Use YYYY-MM-DD"})
		}
	}

	// Parse end date
	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid end_date format. Use YYYY-MM-DD"})
		}
		// Set end date to end of day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	var result struct {
		TotalCustomers         int     `json:"total_customers"`
		RepeatCustomers        int     `json:"repeat_customers"`
		RepeatCustomersRevenue float64 `json:"repeat_customers_revenue"`
	}

	var totalShipping float64
	var totalDeliveryFee float64

	// Total Shipping Fee
	shippingQuery := database.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(shipping_fee), 0)").
		Where("order_status = ?", "delivered")

	// Apply date filters if provided
	if startDateStr != "" {
		shippingQuery = shippingQuery.Where("delivered_at >= ?", startDate)
	}
	if endDateStr != "" {
		shippingQuery = shippingQuery.Where("delivered_at <= ?", endDate)
	}

	err = shippingQuery.Scan(&totalShipping).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch shipping fee"})
	}

	// Total Delivery Partner Fee
	deliveryQuery := database.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(delivery_partner_fee), 0)").
		Where("order_status = ?", "delivered")

	// Apply date filters if provided
	if startDateStr != "" {
		deliveryQuery = deliveryQuery.Where("delivered_at >= ?", startDate)
	}
	if endDateStr != "" {
		deliveryQuery = deliveryQuery.Where("delivered_at <= ?", endDate)
	}

	err = deliveryQuery.Scan(&totalDeliveryFee).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch delivery partner fee"})
	}

	// Repeat customer stats via raw SQL
	query := `
	SELECT
	COUNT(DISTINCT oi.user_id) AS total_customers,
	COUNT(DISTINCT CASE WHEN user_order_count > 1 THEN oi.user_id END) AS repeat_customers,
	COALESCE(SUM(CASE WHEN user_order_count > 1 THEN oi.sub_total ELSE 0 END), 0) AS repeat_customers_revenue,
	COUNT(DISTINCT CASE WHEN user_order_count = 1 THEN oi.user_id END) AS one_time_customers,
	COALESCE(SUM(CASE WHEN user_order_count = 1 THEN oi.sub_total ELSE 0 END), 0) AS one_time_customers_revenue
	FROM (
		SELECT oi_inner.user_id, COUNT(DISTINCT oi_inner.order_id) AS user_order_count
		FROM order_items oi_inner
		JOIN orders o_inner ON oi_inner.order_id = o_inner.order_id
		WHERE o_inner.order_status = 'delivered'`

	// Add date filters to the subquery
	if startDateStr != "" {
		query += ` AND o_inner.delivered_at >= ?`
	}
	if endDateStr != "" {
		query += ` AND o_inner.delivered_at <= ?`
	}

	query += `
		GROUP BY oi_inner.user_id
	) AS user_orders
	JOIN order_items oi ON oi.user_id = user_orders.user_id
	JOIN orders o ON oi.order_id = o.order_id
	WHERE o.order_status = 'delivered'`

	// Add date filters to the main query
	if startDateStr != "" {
		query += ` AND o.delivered_at >= ?`
	}
	if endDateStr != "" {
		query += ` AND o.delivered_at <= ?`
	}

	// Prepare query parameters
	var queryParams []interface{}
	if startDateStr != "" {
		queryParams = append(queryParams, startDate)
	}
	if endDateStr != "" {
		queryParams = append(queryParams, endDate)
	}
	if startDateStr != "" {
		queryParams = append(queryParams, startDate)
	}
	if endDateStr != "" {
		queryParams = append(queryParams, endDate)
	}

	err = database.DB.Raw(query, queryParams...).Scan(&result).Error
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch customer stats"})
	}

	return c.JSON(fiber.Map{
		"result":                     result,
		"total_shipping_revenue":     totalShipping,
		"total_delivery_partner_fee": totalDeliveryFee,
	})
}
