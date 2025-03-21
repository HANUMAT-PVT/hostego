package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"
	"encoding/json"
	"strconv"

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
	ActualShippingFee    float64 `json:"actual_shipping_fee"`
}

// Move struct definition outside the function
type requestCreateOrder struct {
	AddressId       string `json:"address_id"`
	CookingRequests string `json:"cooking_requests"`
}

func CreateNewOrder(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	freeDelivery := false
	var cartItems []models.CartItem
	var order models.Order

	var order_items []models.Order

	var requestOrder requestCreateOrder // Use the struct here
	if err := c.Bind().JSON(&requestOrder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if err := database.DB.Where("user_id=?", user_id).Find(&order_items).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if len(order_items) < 1 {
		freeDelivery = true
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
	totalCharges := CalculateFinalOrderValue(cartItems, freeDelivery)
	order.PlatformFee = totalCharges.PlatformFee
	order.ShippingFee = totalCharges.ShippingFee
	order.FinalOrderValue = totalCharges.FinalOrderValue
	order.DeliveryPartnerFee = totalCharges.DeliveryPartnerShare
	order.OrderStatus = "pending"
	order.AddressID = requestOrder.AddressId
	order.FreeDelivery = freeDelivery
	order.CookingRequests = requestOrder.CookingRequests

	if err := database.DB.Preload("User").Create(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

// CalculateFinalOrderValue calculates the total order cost including charges
func CalculateFinalOrderValue(cartItems []models.CartItem, freeDelivery bool) FinalOrderValueType {
	totalItemSubTotal := 0.0

	// Calculate the subtotal of all cart items
	for _, item := range cartItems {
		totalItemSubTotal = item.SubTotal + totalItemSubTotal
	}

	// Calculate charges
	var shippingFee float64
	if totalItemSubTotal <= 150.0 {
		shippingFee = 15.0
	} else {
		charge := totalItemSubTotal * 0.10
		if charge > 30.0 {
			shippingFee = 30.0
		} else {
			shippingFee = math.Round(charge*100) / 100 // Round to 2 decimal places
		}
	}

	// Platform fee (fixed at ₹1)
	platformFee := 1.0
	// Distribution of charge (80% for delivery partner, 20% for company)
	deliveryPartnerShare := math.Round((shippingFee*0.8)*100) / 100

	shippingFee += platformFee
	actualShippingFee := shippingFee
	if freeDelivery {
		shippingFee = 0
	}
	// Final order value including charges
	finalOrderValue := math.Round((totalItemSubTotal+shippingFee)*100) / 100
	if totalItemSubTotal == 0 {
		finalOrderValue = 0
		shippingFee = 0
		deliveryPartnerShare = 0
		platformFee = 0
		actualShippingFee = 0
	}
	return FinalOrderValueType{
		SubTotal:             math.Round(totalItemSubTotal*100) / 100,
		ShippingFee:          shippingFee,
		PlatformFee:          platformFee,
		DeliveryPartnerShare: deliveryPartnerShare,
		FinalOrderValue:      finalOrderValue,
		ActualShippingFee:    actualShippingFee,
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
	if delivery_partner.AvailabilityStatus == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Delivery partner is not available"})
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

	var delivery_partner models.DeliveryPartner
	if err := database.DB.Where("delivery_partner_id = ?", existingOrder.DeliveryPartnerId).First(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Delivery partner not found"})
	}

	var updateData struct {
		OrderStatus       models.OrderStatusType `json:"order_status"`
		DeliveryPartnerId string                 `json:"delivery_partner_id"`
	}
	if err := c.Bind().JSON(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Update only the status
	existingOrder.OrderStatus = updateData.OrderStatus
	if updateData.DeliveryPartnerId != "" {
		existingOrder.DeliveryPartnerId = updateData.DeliveryPartnerId
	}
	if updateData.OrderStatus == models.DeliveredOrderStatus {
		existingOrder.DeliveredAt = time.Now()
		AddEarningsToDeliveryPartnerWallet(existingOrder)

	}

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
	queryPage := c.Query("page", "1")
	queryLimit := c.Query("limit", "10")
	page, err := strconv.Atoi(queryPage)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(queryLimit)
	if err != nil || limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	if user_id == "" {

	}
	var orders []models.Order
	if err := database.DB.Preload("User").Preload("PaymentTransaction").Preload("Address").Where("user_id=?", user_id).Order("created_at desc").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
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
	if err := dbQuery.Preload("User").Preload("PaymentTransaction").Preload("Address").Order("created_at asc").Find(&orders).Error; err != nil {
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
		if status == "active" {
			statuses := []string{models.AssignedOrderStatus, models.ReachedOrderStatus, models.OnTheWayOrderStatus, models.PackedOrderStatus, models.PickedOrderStatus, models.ReachedDoorStatus, models.CookingOrderStatus}
			dbQuery = dbQuery.Where("order_status IN ?", statuses)
		} else {
			dbQuery = dbQuery.Where("order_status = ?", status)
		}
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

func FetchAllOrderItemsAccordingToProducts(c fiber.Ctx) error {
	type ProductStats struct {
		ProductId     string  `json:"product_id"`
		ProductName   string  `json:"product_name"`
		ProductImg    string  `json:"product_img_url"`
		OrderCount    int     `json:"order_count"`
		TotalQuantity int     `json:"total_quantity"`
		TotalRevenue  float64 `json:"total_revenue"`
		StockQuantity int     `json:"stock_quantity"`
		Availability  int     `json:"availability"`
		Description   string  `json:"description"`
		CurrentPrice  float64 `json:"current_price"`
		// Date-based stats
		LastDayRevenue   float64 `json:"last_day_revenue"`
		LastWeekRevenue  float64 `json:"last_week_revenue"`
		LastMonthRevenue float64 `json:"last_month_revenue"`
		LastDayOrders    int     `json:"last_day_orders"`
		LastWeekOrders   int     `json:"last_week_orders"`
		LastMonthOrders  int     `json:"last_month_orders"`
	}

	type OverallStats struct {
		TotalRevenue     float64 `json:"total_revenue"`
		TotalOrders      int     `json:"total_orders"`
		LastDayRevenue   float64 `json:"last_day_revenue"`
		LastWeekRevenue  float64 `json:"last_week_revenue"`
		LastMonthRevenue float64 `json:"last_month_revenue"`
		LastDayOrders    int     `json:"last_day_orders"`
		LastWeekOrders   int     `json:"last_week_orders"`
		LastMonthOrders  int     `json:"last_month_orders"`
	}

	// Get date range filters from query params
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Get product-wise stats
	var stats []ProductStats
	query := database.DB.
		Select(`
			products.product_id,
			products.product_name,
			products.product_img as product_img,
			products.description,
			products.stock_quantity,
			products.availability,
			products.food_price as current_price,
			COUNT(DISTINCT CASE WHEN orders.order_id IS NOT NULL THEN order_items.order_id END) as order_count,
			COALESCE(SUM(order_items.quantity), 0) as total_quantity,
			COALESCE(SUM(order_items.sub_total), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN orders.created_at >= NOW() - INTERVAL '1 day' THEN order_items.sub_total ELSE 0 END), 0) as last_day_revenue,
			COALESCE(SUM(CASE WHEN orders.created_at >= NOW() - INTERVAL '7 day' THEN order_items.sub_total ELSE 0 END), 0) as last_week_revenue,
			COALESCE(SUM(CASE WHEN orders.created_at >= NOW() - INTERVAL '30 day' THEN order_items.sub_total ELSE 0 END), 0) as last_month_revenue,
			COUNT(DISTINCT CASE WHEN orders.created_at >= NOW() - INTERVAL '1 day' THEN order_items.order_id END) as last_day_orders,
			COUNT(DISTINCT CASE WHEN orders.created_at >= NOW() - INTERVAL '7 day' THEN order_items.order_id END) as last_week_orders,
			COUNT(DISTINCT CASE WHEN orders.created_at >= NOW() - INTERVAL '30 day' THEN order_items.order_id END) as last_month_orders
		`).
		Table("products").
		Joins("LEFT JOIN order_items ON order_items.product_id::text = products.product_id").
		Joins("LEFT JOIN orders ON orders.order_id = order_items.order_id")

	if startDate != "" && endDate != "" {
		startDateTime := startDate + " 00:00:00"
		endDateTime := endDate + " 23:59:59"
		query = query.Where("orders.created_at BETWEEN ? AND ? OR orders.created_at IS NULL", startDateTime, endDateTime)
	}

	err := query.
		Group("products.product_id").
		Order("products.product_name").
		Scan(&stats).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to fetch product statistics",
			"details": err.Error(),
		})
	}

	// Get overall stats with a separate query
	var overall OverallStats
	err = database.DB.Model(&models.OrderItem{}).
		Select(`
			COALESCE(SUM(order_items.sub_total), 0) as total_revenue,
			COUNT(DISTINCT order_items.order_id) as total_orders,
			COALESCE(SUM(CASE WHEN o.created_at >= NOW() - INTERVAL '1 day' THEN order_items.sub_total ELSE 0 END), 0) as last_day_revenue,
			COALESCE(SUM(CASE WHEN o.created_at >= NOW() - INTERVAL '7 day' THEN order_items.sub_total ELSE 0 END), 0) as last_week_revenue,
			COALESCE(SUM(CASE WHEN o.created_at >= NOW() - INTERVAL '30 day' THEN order_items.sub_total ELSE 0 END), 0) as last_month_revenue,
			COUNT(DISTINCT CASE WHEN o.created_at >= NOW() - INTERVAL '1 day' THEN order_items.order_id END) as last_day_orders,
			COUNT(DISTINCT CASE WHEN o.created_at >= NOW() - INTERVAL '7 day' THEN order_items.order_id END) as last_week_orders,
			COUNT(DISTINCT CASE WHEN o.created_at >= NOW() - INTERVAL '30 day' THEN order_items.order_id END) as last_month_orders
		`).
		Joins("LEFT JOIN orders o ON o.order_id = order_items.order_id").
		Scan(&overall).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to fetch overall statistics",
			"details": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product_stats": stats,
		"overall_stats": overall,
	})
}
