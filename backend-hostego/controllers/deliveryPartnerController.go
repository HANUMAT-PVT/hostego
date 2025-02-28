package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"
	"math"

	"github.com/gofiber/fiber/v3"
)

func CreateNewDeliveryPartner(c fiber.Ctx) error {
	var delivery_partner models.DeliveryPartner
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if err := database.DB.Where("user_id=?", user_id).Find(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}
	delivery_partner.UserId = user_id
	if err := c.Bind().JSON(&delivery_partner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}
	if err := database.DB.Preload("User").Create(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"err": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"delivery_partner": delivery_partner})
}

func UpdateDeliveryPartner(c fiber.Ctx) error {
	var delivery_partner models.DeliveryPartner
	delivery_partner_id := c.Params("id")
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}

	// First find the existing delivery partner
	if err := database.DB.Where("delivery_partner_id = ? AND user_id = ?", delivery_partner_id, user_id).First(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Delivery partner not found"})
	}

	// Bind the updated data
	if err := c.Bind().JSON(&delivery_partner); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Save the updates
	if err := database.DB.Save(&delivery_partner).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"delivery_partner": delivery_partner,
		"message":          "Delivery Partner Updated successfully",
	})
}

func FetchDeliveryPartnerByUserId(c fiber.Ctx) error {

	user_id, err := middlewares.VerifyUserAuthCookie(c)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error(), "message": "You are not Authenticated !"})
	}
	if user_id != "" {
	}

	var delivery_partner models.DeliveryPartner

	if err := database.DB.Preload("User").First(&delivery_partner, "user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Delivery partner not found"})
	}
	return c.Status(fiber.StatusOK).JSON(delivery_partner)
}

func FetchAllDeliveryPartners(c fiber.Ctx) error {
	var delivery_partners []models.DeliveryPartner
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You are not Authenticated !"})
	}
	dbQuery := database.DB

	availability := c.Query("availability")
	account_status := c.Query("account_status")
	verification_status := c.Query("verification_status")

	if availability != "" {
		dbQuery = dbQuery.Where("availability_status=?", availability)
	}
	if account_status != "" {
		dbQuery = dbQuery.Where("account_status=?", account_status)
	}
	if verification_status != "" {
		dbQuery = dbQuery.Where("verification_status=?", verification_status)
	}

	if err := dbQuery.Preload("User").Find(&delivery_partners).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(delivery_partners)
}

type DeliveryEarningStats struct {
	Date            string         `json:"date"`
	OrdersDelivered int            `json:"orders_delivered"`
	TotalEarning    float64        `json:"total_earning"`
	Orders          []models.Order `json:"orders"`
}

func FetchDeliveryPartnerEarnings(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "You are not Authenticated!"})
	}

	delivery_partner_id := c.Params("id")

	// Optional date range filters
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	dbQuery := database.DB.
		Preload("User").
		Preload("PaymentTransaction").
		Where("delivery_partner_id = ?", delivery_partner_id).
		Where("order_status = ?", models.DeliveredOrderStatus)

	// Add date range filter if provided
	if startDate != "" && endDate != "" {
		dbQuery = dbQuery.Where("delivered_at BETWEEN ? AND ?", startDate, endDate)
	}

	var orders []models.Order
	if err := dbQuery.Order("delivered_at desc").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Group orders by date and calculate statistics
	dailyStats := make(map[string]*DeliveryEarningStats)
	var totalEarnings float64
	totalOrders := len(orders)

	for _, order := range orders {
		dateStr := order.DeliveredAt.Format("2006-01-02") // YYYY-MM-DD format

		// Initialize stats for this date if not exists
		if _, exists := dailyStats[dateStr]; !exists {
			dailyStats[dateStr] = &DeliveryEarningStats{
				Date:            dateStr,
				OrdersDelivered: 0,
				TotalEarning:    0,
				Orders:          []models.Order{},
			}
		}

		// Update daily statistics
		stats := dailyStats[dateStr]
		stats.OrdersDelivered++
		stats.TotalEarning += order.DeliveryPartnerFee
		stats.Orders = append(stats.Orders, order)

		totalEarnings += order.DeliveryPartnerFee
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"daily_earnings": orders,
		"summary": fiber.Map{
			"total_earnings": math.Round(totalEarnings*100) / 100,
			"total_orders":   totalOrders,
		},
	})
}
