package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"encoding/json"
	"strconv"

	"math"

	"time"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type FinalOrderValueType struct {
	SubTotal             float64 `json:"subtotal"`
	ShippingFee          float64 `json:"shipping_fee"`
	PlatformFee          float64 `json:"platform_fee"`
	DeliveryPartnerShare float64 `json:"delivery_partner_fee"`
	FinalOrderValue      float64 `json:"final_order_value"`
	ActualShippingFee    float64 `json:"actual_shipping_fee"`
	RainExtraFee         float64 `json:"rain_extra_fee"`
}

// Move struct definition outside the function
type requestCreateOrder struct {
	AddressId       int    `json:"address_id"`
	CookingRequests string `json:"cooking_requests"`
}

func CreateNewOrder(c fiber.Ctx) error {
	user_id := c.Locals("user_id").(int)
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

	if requestOrder.AddressId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Address ID is required"})
	}

	if err := database.DB.Preload("ProductItem.Shop").Where("user_id=?", user_id).Find(&cartItems).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	// delete the cart item if the shop is closed or not
	for _, cartItem := range cartItems {
		if cartItem.ProductItem.Shop.ShopStatus == 0 {
			database.DB.Where("cart_item_id = ? AND user_id = ?", cartItem.CartItemId, user_id).Delete(&cartItem)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Shop is closed!"})
		}
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
	rainExtraCharge := 10.0
	// Distribution of charge (80% for delivery partner, 20% for company)
	deliveryPartnerShare := math.Round((shippingFee*0.8)*100) / 100
	// deliveryPartnerShare += rainExtraCharge * 0.7
	shippingFee = 0
	shippingFee += platformFee

	if totalItemSubTotal <= 150.0 {
		shippingFee += 150 * 0.15
	} else {
		charge := totalItemSubTotal * 0.15
		if charge > 39.0 {
			shippingFee = 39.0
		} else {
			shippingFee += math.Round(charge*100) / 100 // Round to 2 decimal places
		}
	}
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
		rainExtraCharge = 10.0
	}
	if rainExtraCharge != 0 {
		shippingFee += rainExtraCharge
		finalOrderValue += rainExtraCharge;
	}
	return FinalOrderValueType{
		SubTotal:             math.Round(totalItemSubTotal*100) / 100,
		ShippingFee:          shippingFee,
		PlatformFee:          platformFee,
		DeliveryPartnerShare: deliveryPartnerShare,
		FinalOrderValue:      finalOrderValue,
		ActualShippingFee:    actualShippingFee,
		RainExtraFee:         10.0,
	}
}

func MarkOrderAsDelivered(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {

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
	user_id := c.Locals("user_id")

	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	order_id := c.Params("id")
	var order models.Order

	if err := database.DB.Preload("User").Preload("PaymentTransaction").Preload("Address").Where("order_id=?", order_id).First(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusOK).JSON(order)
}

func AssignOrderToDeliveryPartner(c fiber.Ctx) error {
	user_id := c.Locals("user_id")

	type requestAssignOrder struct {
		DeliveryPartnerId int `json:"delivery_partner_id"`
		OrderId           int `json:"order_id"`
	}
	var request_assign requestAssignOrder

	var order models.Order

	var delivery_partner models.DeliveryPartner

	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
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
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	order_id := c.Params("id")

	var existingOrder models.Order
	if err := database.DB.Where("order_id = ?", order_id).First(&existingOrder).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	// var delivery_partner models.DeliveryPartner
	// if err := database.DB.Where("delivery_partner_id = ?", existingOrder.DeliveryPartnerId).First(&delivery_partner).Error; err != nil {
	// 	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Delivery partner not found"})
	// }

	var updateData struct {
		OrderStatus       models.OrderStatusType `json:"order_status"`
		DeliveryPartnerId int                    `json:"delivery_partner_id"`
	}
	if err := c.Bind().JSON(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Update only the status
	if updateData.OrderStatus == models.CanceledOrderStatus {
		if existingOrder.OrderStatus == models.DeliveredOrderStatus {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Order can't be cancelled its delivered already !",
				"order":   existingOrder,
			})
		}
	}
	if updateData.OrderStatus == models.DeliveredOrderStatus {
		// delivered status only when order is not delivered yet
		if existingOrder.OrderStatus != models.DeliveredOrderStatus {

			existingOrder.DeliveredAt = time.Now()
			AddEarningsToDeliveryPartnerWallet(existingOrder)
		}

	}
	existingOrder.OrderStatus = updateData.OrderStatus
	if updateData.DeliveryPartnerId != 0 {
		existingOrder.DeliveryPartnerId = updateData.DeliveryPartnerId
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
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
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

	if user_id == 0 {

	}
	var orders []models.Order
	if err := database.DB.Preload("User").Preload("PaymentTransaction").Preload("Address").Where("user_id=?", user_id).Order("created_at desc").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(orders)
}

func FetchAllOrders(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
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

	dbQuery := database.DB

	searchQuery := c.Query("search")

	filter := c.Query("filter")

	if searchQuery != "" {
		dbQuery = dbQuery.Where(
			"order_id = ? OR user_id = ? OR mobile_number = ?",
			searchQuery, searchQuery, searchQuery,
		)
	}

	if filter != "" {
		dbQuery = dbQuery.Where("order_status = ?", filter)
	}

	var orders []models.Order
	if err := dbQuery.Preload("User").Preload("PaymentTransaction").Preload("Address").Order("created_at desc").Limit(limit).Offset(offset).Find(&orders).Order("created_at desc").Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}

func FetchAllOrdersByDeliveryPartner(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
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
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	type ProductStats struct {
		ProductId        string  `json:"product_id"`
		ProductName      string  `json:"product_name"`
		ProductImg       string  `json:"product_img_url"`
		OrderCount       int     `json:"order_count"`
		TotalQuantity    int     `json:"total_quantity"`
		TotalRevenue     float64 `json:"total_revenue"`
		StockQuantity    int     `json:"stock_quantity"`
		Availability     int     `json:"availability"`
		Description      string  `json:"description"`
		CurrentPrice     float64 `json:"current_price"`
		LastDayRevenue   float64 `json:"last_day_revenue"`
		LastWeekRevenue  float64 `json:"last_week_revenue"`
		LastMonthRevenue float64 `json:"last_month_revenue"`
		LastDayOrders    int     `json:"last_day_orders"`
		LastWeekOrders   int     `json:"last_week_orders"`
		LastMonthOrders  int     `json:"last_month_orders"`
		ShopName         string  `json:"shop_name"`
	}

	// Get date range filters from query params
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	// Base query to get product stats
	query := database.DB.Model(&models.Product{}).
		Select(`
			products.product_id,
			products.product_name,
			products.product_img,
			products.description,
			products.stock_quantity,
			products.availability,
			shops.shop_name as shop_name,
			products.food_price as current_price,
			COUNT(DISTINCT order_items.order_id) as order_count,
			COALESCE(SUM(order_items.quantity), 0) as total_quantity,
			COALESCE(SUM(order_items.sub_total), 0) as total_revenue,
			COALESCE(SUM(CASE WHEN orders.created_at >= NOW() - INTERVAL '1 day' THEN order_items.sub_total ELSE 0 END), 0) as last_day_revenue,
			COALESCE(SUM(CASE WHEN orders.created_at >= NOW() - INTERVAL '7 day' THEN order_items.sub_total ELSE 0 END), 0) as last_week_revenue,
			COALESCE(SUM(CASE WHEN orders.created_at >= NOW() - INTERVAL '30 day' THEN order_items.sub_total ELSE 0 END), 0) as last_month_revenue,
			COUNT(DISTINCT CASE WHEN orders.created_at >= NOW() - INTERVAL '1 day' THEN order_items.order_id END) as last_day_orders,
			COUNT(DISTINCT CASE WHEN orders.created_at >= NOW() - INTERVAL '7 day' THEN order_items.order_id END) as last_week_orders,
			COUNT(DISTINCT CASE WHEN orders.created_at >= NOW() - INTERVAL '30 day' THEN order_items.order_id END) as last_month_orders
		`).
		Joins("LEFT JOIN shops ON shops.shop_id = products.shop_id").
		Joins("LEFT JOIN order_items ON order_items.product_id = products.product_id").
		Joins("LEFT JOIN orders ON orders.order_id = order_items.order_id")

	if startDate != "" && endDate != "" {
		query = query.Where("orders.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}

	var stats = []ProductStats{}
	err := query.Group(`
		products.product_id, 
		products.product_name, 
		products.product_img, 
		products.description, 
		products.stock_quantity, 
		products.availability, 
		shops.shop_name,
		products.food_price
	`).
		Order("products.product_name").
		Scan(&stats).Error

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Failed to fetch product statistics",
			"details": err.Error(),
		})
	}

	// Calculate overall stats
	var overallStats struct {
		TotalRevenue     float64 `json:"total_revenue"`
		TotalOrders      int     `json:"total_orders"`
		LastDayRevenue   float64 `json:"last_day_revenue"`
		LastWeekRevenue  float64 `json:"last_week_revenue"`
		LastMonthRevenue float64 `json:"last_month_revenue"`
		LastDayOrders    int     `json:"last_day_orders"`
		LastWeekOrders   int     `json:"last_week_orders"`
		LastMonthOrders  int     `json:"last_month_orders"`
	}

	for _, stat := range stats {
		overallStats.TotalRevenue += stat.TotalRevenue
		overallStats.TotalOrders += stat.OrderCount
		overallStats.LastDayRevenue += stat.LastDayRevenue
		overallStats.LastWeekRevenue += stat.LastWeekRevenue
		overallStats.LastMonthRevenue += stat.LastMonthRevenue
		overallStats.LastDayOrders += stat.LastDayOrders
		overallStats.LastWeekOrders += stat.LastWeekOrders
		overallStats.LastMonthOrders += stat.LastMonthOrders
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"product_stats": stats,
		"overall_stats": overallStats,
	})
}

func CancelOrder(c fiber.Ctx) error {

	current_user_id := c.Locals("user_id").(int)
	if current_user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	type OrderRequest struct {
		OrderID int `json:"order_id"`
	}

	var request OrderRequest

	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if current_user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	var order models.Order

	var delivery_partner models.DeliveryPartner

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := database.DB.Where("order_id = ?", request.OrderID).First(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if order.DeliveryPartnerId != 0 {
		if err := tx.Where("delivery_partner_id = ?", order.DeliveryPartnerId).First(&delivery_partner).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	order.OrderStatus = models.OrderStatusType(models.CanceledOrderStatus)
	order.DeliveryPartnerId = 0
	order.DeliveryPartner = nil

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var orderItems []models.CartItem
	if err := json.Unmarshal(order.OrderItems, &orderItems); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse order items"})
	}

	for _, item := range orderItems {
		if err := tx.Model(&models.Product{}).Where("product_id = ?", item.ProductId).
			Update("stock_quantity", gorm.Expr("stock_quantity + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := tx.Where("order_id = ?", order.OrderId).
			Delete(&models.OrderItem{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order  Cancelled/without refund"})

}
