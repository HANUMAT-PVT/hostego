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

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

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

	// Build base query
	revenueQuery := database.DB.
		Table("order_items").
		Select(`
			COALESCE(SUM(order_items.actual_sub_total), 0) AS total_revenue,
			COUNT(DISTINCT order_items.order_id) AS total_orders
		`).
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("products.shop_id = ? AND orders.order_status = ?", shopID, "delivered")

	// Apply date filter if provided
	if startDate != "" && endDate != "" {
		revenueQuery = revenueQuery.Where("orders.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	if err := revenueQuery.Scan(&revenueStats).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate revenue/orders"})
	}

	if revenueStats.TotalOrders > 0 {
		revenueStats.AverageOrderValue = revenueStats.TotalRevenue / float64(revenueStats.TotalOrders)
	}

	// Average product rating
	ratingQuery := database.DB.
		Table("products").
		Select("AVG(products.average_rating) AS average_rating").
		Joins("JOIN order_items ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("products.shop_id = ? AND orders.order_status = ?", shopID, "delivered")

	if startDate != "" && endDate != "" {
		ratingQuery = ratingQuery.Where("orders.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	if err := ratingQuery.Scan(&averageStats).Error; err != nil {
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

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

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

	query := database.DB.
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
		Where("products.shop_id = ? AND orders.order_status = ?", shopID, "delivered")

	if startDate != "" && endDate != "" {
		query = query.Where("orders.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	query.Group("products.product_id, products.product_name, products.food_price, products.selling_price, products.average_rating").
		Group("products.product_id, products.product_name, products.food_price, products.selling_price, products.average_rating").
		Order("quantity_sold DESC").
		Limit(limit).
		Scan(&topProducts)

	if query.Error != nil {
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

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	type Result struct {
		TotalOrderCount     int64   `json:"total_orders"`
		DeliveredOrderCount int64   `json:"delivered_orders"`
		CancelledOrderCount int64   `json:"cancelled_orders"`
		TotalRevenue        float64 `json:"total_revenue"`
		AverageOrderValue   float64 `json:"average_order_value"`
	}

	var result Result

	// Total Orders by shop_id (from orders)
	query := database.DB.
		Model(&models.Order{}).
		Where("shop_id = ?", shopID)

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	if err := query.Count(&result.TotalOrderCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching total orders"})
	}

	// Delivered Orders
	query = database.DB.
		Model(&models.Order{}).
		Where("shop_id = ? AND order_status = ?", shopID, "delivered")

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	if err := query.Count(&result.DeliveredOrderCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching delivered orders"})
	}

	// Cancelled Orders
	query = database.DB.
		Model(&models.Order{}).
		Where("shop_id = ? AND order_status = ?", shopID, "cancelled")

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	if err := query.Count(&result.CancelledOrderCount).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error fetching cancelled orders"})
	}

	// Total Revenue from order_items for delivered orders
	revenueQuery := database.DB.
		Table("order_items").
		Select("COALESCE(SUM(order_items.actual_sub_total), 0)").
		Joins("JOIN orders ON orders.order_id = order_items.order_id").
		Where("orders.shop_id = ? AND orders.order_status = ?", shopID, "delivered")

	if startDate != "" && endDate != "" {
		revenueQuery = revenueQuery.Where("orders.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err = revenueQuery.Scan(&result.TotalRevenue).Error
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

	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	db := database.DB // your GORM DB instance

	var totalCustomers int64
	totalCustomersQuery := db.
		Table("order_items").
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("products.shop_id = ? AND orders.order_status = ?", shopID, "delivered")

	if startDate != "" && endDate != "" {
		totalCustomersQuery = totalCustomersQuery.Where("orders.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err := totalCustomersQuery.Select("COUNT(DISTINCT order_items.user_id)").Scan(&totalCustomers).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to count customers"})
	}

	var repeatCustomers int64
	repeatCustomersQuery := db.
		Table("order_items").
		Joins("JOIN products ON order_items.product_id = products.product_id").
		Joins("JOIN orders ON order_items.order_id = orders.order_id").
		Where("products.shop_id = ? AND orders.order_status = ?", shopID, "delivered")

	if startDate != "" && endDate != "" {
		repeatCustomersQuery = repeatCustomersQuery.Where("orders.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	err = repeatCustomersQuery.Select("COUNT(*)").
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
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// 1. Average Preparation Time
	var prepDurations []time.Duration
	var orders []models.Order

	prepQuery := database.DB.Where("shop_id = ? AND order_status = ? AND restaurant_responded_at IS NOT NULL AND actual_ready_at IS NOT NULL", shopID, "delivered")

	if startDate != "" && endDate != "" {
		prepQuery = prepQuery.Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	if err := prepQuery.Find(&orders).Error; err != nil {
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

	peakHoursQuery := database.DB.Table("orders").
		Select("EXTRACT(HOUR FROM created_at) AS hour, COUNT(*) AS order_count").
		Where("shop_id = ? AND order_status = ?", shopID, "delivered")

	if startDate != "" && endDate != "" {
		peakHoursQuery = peakHoursQuery.Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	if err := peakHoursQuery.Group("hour").
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

func GetRestaurantRevenueAnalytics(c *fiber.Ctx) error {
	shopID := c.Params("shop_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	type RevenueAnalytics struct {
		TotalRevenue float64 `json:"total_revenue"`
		TotalPending float64 `json:"total_pending"`
	}
	var revenueAnalytics RevenueAnalytics

	query := database.DB.
		Table("orders").
		Where("shop_id = ? AND order_status = ?", shopID, "delivered")

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	query.Select("SUM(restaurant_payable_amount) AS total_revenue").
		Scan(&revenueAnalytics)

	if err := query.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch revenue"})
	}

	// Pending revenue query
	pendingQuery := database.DB.
		Table("orders").
		Where("shop_id = ? AND restaurant_paid_at IS NULL", shopID)

	if startDate != "" && endDate != "" {
		pendingQuery = pendingQuery.Where("created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	pendingQuery.Select("SUM(restaurant_payable_amount) AS total_pending").
		Scan(&revenueAnalytics.TotalPending)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Restaurant revenue analytics",
		"total_revenue": revenueAnalytics.TotalRevenue,
		"total_pending": revenueAnalytics.TotalPending,
	})
}
