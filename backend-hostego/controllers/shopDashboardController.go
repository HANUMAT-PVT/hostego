package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetShopDashboardStats(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	if shopIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "shop_id is required"})
	}
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid shop_id"})
	}

	timeRange := c.Query("range", "day") // default range
	var startTime time.Time
	now := time.Now()

	switch timeRange {
	case "day":
		startTime = now.Add(-24 * time.Hour)
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid time range"})
	}

	type AverageStats struct {
		AverageRating float64 `json:"average_product_rating"`
	}

	type TotalRevenueStats struct {
		TotalRevenue      float64 `json:"total_revenue"`
		TotalOrders       int     `json:"total_orders"`
		AverageOrderValue float64 `json:"average_order_value"`
	}

	var revenueStats TotalRevenueStats
	var averageStats AverageStats

	// 1. Total revenue & order count
	err = database.DB.
		Table("order_items").
		Select(`
				COALESCE(SUM(order_items.actual_sub_total), 0) AS total_revenue,
				COUNT(DISTINCT order_items.order_id) AS total_orders
			`).
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("products.shop_id = ? AND orders.order_status = ? AND orders.created_at >= ?", shopID, "delivered", startTime).
		Scan(&revenueStats).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate revenue/orders"})
	}

	if revenueStats.TotalOrders > 0 {
		revenueStats.AverageOrderValue = revenueStats.TotalRevenue / float64(revenueStats.TotalOrders)
	}

	// 2. Average rating for delivered products
	err = database.DB.
		Table("products").
		Select("AVG(products.average_rating) AS average_rating").
		Joins("JOIN order_items ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("products.shop_id = ? AND orders.order_status = ? AND orders.created_at >= ?", shopID, "delivered", startTime).
		Scan(&averageStats).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch average product rating"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Shop dashboard",
		"stats": fiber.Map{
			"average_product_rating": averageStats.AverageRating,
			"total_revenue":          revenueStats.TotalRevenue,
			"total_orders":           revenueStats.TotalOrders,
			"average_order_value":    revenueStats.AverageOrderValue,
		},
	})
}

func GetTopSellingProducts(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	if shopIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "shop_id is required"})
	}
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid shop_id"})
	}

	rangeStr := c.Query("range", "week") // default to week
	var startTime time.Time
	now := time.Now()

	switch rangeStr {
	case "day":
		startTime = now.Add(-24 * time.Hour)
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid time range"})
	}

	limitStr := c.Query("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	type TopProduct struct {
		ProductID     int     `json:"product_id"`
		ProductName   string  `json:"product_name"`
		QuantitySold  int     `json:"quantity_sold"`
		FoodPrice     float64 `json:"food_price"`
		SellingPrice  float64 `json:"selling_price"`
		AverageRating float64 `json:"average_rating"`
	}

	var topProducts []TopProduct

	err = database.DB.
		Table("order_items").
		Select(`
			products.product_id,
			products.product_name,
			products.food_price,
			products.selling_price,
			products.average_rating,
			SUM(order_items.quantity) AS quantity_sold
		`).
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("products.shop_id = ? AND orders.order_status = ? AND orders.created_at >= ?", shopID, "delivered", startTime).
		Group("products.product_id, products.product_name, products.food_price, products.selling_price, products.average_rating").
		Order("quantity_sold DESC").
		Limit(limit).
		Scan(&topProducts).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch top selling products"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"top_selling_products": topProducts,
	})
}

func GetOrderAnalytics(c *fiber.Ctx) error {
	shopIDStr := c.Params("shop_id")
	if shopIDStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "shop_id is required"})
	}
	shopID, err := strconv.Atoi(shopIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid shop_id"})
	}

	// Time range: day, week, month
	rangeStr := c.Query("range", "week")
	now := time.Now()
	var startTime time.Time

	switch rangeStr {
	case "day":
		startTime = now.AddDate(0, 0, -1)
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid time range"})
	}

	type Result struct {
		TotalOrderCount     int64   `json:"total_orders"`
		DeliveredOrderCount int64   `json:"delivered_orders"`
		CancelledOrderCount int64   `json:"cancelled_orders"`
		TotalRevenue        float64 `json:"total_revenue"`
		AverageOrderValue   float64 `json:"average_order_value"`
	}

	var result Result

	// Total Orders by shop_id (from orders)
	if err := database.DB.
		Model(&models.Order{}).
		Where("shop_id = ? AND created_at >= ?", shopID, startTime).
		Count(&result.TotalOrderCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching total orders"})
	}

	// Delivered Orders
	if err := database.DB.
		Model(&models.Order{}).
		Where("shop_id = ? AND order_status = ? AND created_at >= ?", shopID, "delivered", startTime).
		Count(&result.DeliveredOrderCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching delivered orders"})
	}

	// Cancelled Orders
	if err := database.DB.
		Model(&models.Order{}).
		Where("shop_id = ? AND order_status = ? AND created_at >= ?", shopID, "cancelled", startTime).
		Count(&result.CancelledOrderCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching cancelled orders"})
	}

	// Total Revenue from order_items for delivered orders
	err = database.DB.
		Table("order_items").
		Select("COALESCE(SUM(order_items.actual_sub_total), 0)").
		Joins("JOIN orders ON orders.order_id = order_items.order_id").
		Where("orders.shop_id = ? AND orders.order_status = ? AND orders.created_at >= ?", shopID, "delivered", startTime).
		Scan(&result.TotalRevenue).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching revenue"})
	}

	// Calculate average order value
	if result.DeliveredOrderCount > 0 {
		result.AverageOrderValue = result.TotalRevenue / float64(result.DeliveredOrderCount)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func GetCustomerInsights(c *fiber.Ctx) error {
	shopIDParam := c.Params("shop_id")
	if shopIDParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "shop_id is required",
		})
	}

	var shopID int
	if id, err := strconv.Atoi(shopIDParam); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid shop_id",
		})
	} else {
		shopID = id
	}

	db := database.DB // your GORM DB instance

	var totalCustomers int64
	err := db.
		Table("order_items").
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("products.shop_id = ? AND orders.order_status = ?", shopID, "delivered").
		Select("COUNT(DISTINCT order_items.user_id)").
		Scan(&totalCustomers).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to count customers"})
	}

	var repeatCustomers int64
	err = db.
		Table("order_items").
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("products.shop_id = ? AND orders.order_status = ?", shopID, "delivered").
		Select("COUNT(*)").
		Group("order_items.user_id").
		Having("COUNT(DISTINCT order_items.order_id) > 1").
		Count(&repeatCustomers).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to count repeat customers"})
	}

	return c.JSON(fiber.Map{
		"total_customers":  totalCustomers,
		"repeat_customers": repeatCustomers,
	})
}

func GetRestaurantPerformanceMetrics(c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	rangeStr := c.Query("range", "month")

	var startTime time.Time
	now := time.Now()

	switch rangeStr {
	case "day":
		startTime = now.AddDate(0, 0, -1)
	case "week":
		startTime = now.AddDate(0, 0, -7)
	case "month":
		startTime = now.AddDate(0, -1, 0)
	default:
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid time range"})
	}

	// 1. Average Preparation Time
	var prepDurations []time.Duration
	var orders []models.Order
	if err := database.DB.Where("shop_id = ? AND order_status = ? AND restaurant_responded_at IS NOT NULL AND actual_ready_at IS NOT NULL AND created_at >= ?", shopID, "delivered", startTime).
		Find(&orders).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch orders"})
	}

	for _, o := range orders {
		if !o.RestaurantRespondedAt.IsZero() && !o.ActualReadyAt.IsZero() {
			prepDurations = append(prepDurations, o.ActualReadyAt.Sub(o.RestaurantRespondedAt))
		}
	}

	var avgPrepMinutes float64
	if len(prepDurations) > 0 {
		var total time.Duration
		for _, d := range prepDurations {
			total += d
		}
		avgPrepMinutes = total.Minutes() / float64(len(prepDurations))
	}

	// 2. Customer Satisfaction from Product Ratings
	var avgRating float64
	if err := database.DB.Table("products").Where("shop_id = ?", shopID).
		Select("AVG(average_rating)").Scan(&avgRating).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to calculate customer satisfaction"})
	}

	// 3. Peak Hours (Hourly Order Count)
	type HourStat struct {
		Hour         int
		OrderCount   int
		TimeInterval string
	}
	var hourStats []HourStat
	if err := database.DB.Table("orders").
		Select("EXTRACT(HOUR FROM created_at) AS hour, COUNT(*) AS order_count").
		Where("shop_id = ? AND order_status = ? AND created_at >= ?", shopID, "delivered", startTime).
		Group("hour").
		Order("hour").
		Scan(&hourStats).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to get peak hour stats"})
	}

	// Format peak hours
	peakHours := make([]fiber.Map, 0)
	for _, hs := range hourStats {
		start := fmt.Sprintf("%02d:00", hs.Hour)
		end := fmt.Sprintf("%02d:00", (hs.Hour+1)%24)
		peakHours = append(peakHours, fiber.Map{
			"time_range": fmt.Sprintf("%s-%s", start, end),
			"orders":     hs.OrderCount,
		})
	}

	return c.JSON(fiber.Map{
		"average_prep_time_minutes": fmt.Sprintf("%.2f", avgPrepMinutes),
		"customer_satisfaction":     fmt.Sprintf("%.2f", avgRating),
		"peak_hours":                peakHours,
	})
}
