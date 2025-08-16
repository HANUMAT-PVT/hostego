package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateRestaurantPayout(shopID int) (*models.RestaurantPayout, error) {
	// 0. If a payout is already pending/processing for this shop, return it instead of creating a new one
	var existing models.RestaurantPayout
	if err := database.DB.Where("shop_id = ? AND status IN (?)", shopID, []string{"pending", "processing"}).Order("created_at DESC").First(&existing).Error; err == nil {
		return &existing, nil
	}

	var unpaidOrders []models.Order

	// 1. Fetch unpaid delivered orders for this shop that are not yet linked to any payout
	if err := database.DB.
		Where("shop_id = ? AND order_status = ? AND restaurant_paid_at IS NULL AND (restaurant_payout_id IS NULL OR restaurant_payout_id = '')", shopID, models.DeliveredOrderStatus).
		Find(&unpaidOrders).Error; err != nil {
		return nil, err
	}

	if len(unpaidOrders) == 0 {
		return nil, fmt.Errorf("no unpaid delivered orders found for shop %d", shopID)
	}

	// 2. Calculate SUM in SQL for the same set of orders
	var total float64
	if err := database.DB.Model(&models.Order{}).
		Where("shop_id = ? AND order_status = ? AND restaurant_paid_at IS NULL AND (restaurant_payout_id IS NULL OR restaurant_payout_id = '')", shopID, models.DeliveredOrderStatus).
		Select("SUM(restaurant_payable_amount)").Scan(&total).Error; err != nil {
		return nil, err
	}

	// 3. Create payout record
	payout := models.RestaurantPayout{
		ShopID:      shopID,
		TotalAmount: total,
		Status:      "pending",
	}

	if err := database.DB.Create(&payout).Error; err != nil {
		return nil, err
	}

	// 4. Link selected orders to the newly created payout to prevent duplication on subsequent runs
	payoutIDStr := strconv.Itoa(payout.PayoutID)
	if err := database.DB.Model(&models.Order{}).
		Where("shop_id = ? AND order_status = ? AND restaurant_paid_at IS NULL AND (restaurant_payout_id IS NULL OR restaurant_payout_id = '')", shopID, models.DeliveredOrderStatus).
		Update("restaurant_payout_id", payoutIDStr).Error; err != nil {
		// Best effort rollback to avoid orphan payout
		_ = database.DB.Delete(&payout).Error
		return nil, err
	}

	return &payout, nil
}

func InitiateRestaurantPayout(c *fiber.Ctx) error {

	var shops []models.Shop
	if err := database.DB.Find(&shops).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Shop not found"})
	}

	for _, shop := range shops {
		_, err := CreateRestaurantPayout(shop.ShopId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to initiate restaurant payout"})
		}

	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Restaurant payout initiated for all shops"})
}

func VerifyRestaurantPayout(c *fiber.Ctx) error {

	payoutID := c.Params("payout_id")

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var request struct {
		PayoutID      string `json:"payout_id"`
		PaymentRefID  string `json:"payment_ref_id"`
		PaymentMethod string `json:"payment_method"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var payout models.RestaurantPayout
	if err := database.DB.First(&payout, "payout_id = ?", payoutID).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payout not found"})
	}

	if payout.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payout is not pending"})
	}

	now := time.Now()
	payout.PaidAt = &now
	payout.PaymentRefID = request.PaymentRefID
	payout.PaymentMethod = request.PaymentMethod
	payout.Status = "paid"

	// update the orders with the payout id
	if err := tx.Model(&models.Order{}).
		Where("restaurant_payout_id = ?", payoutID).
		Update("restaurant_payout_id", payoutID).
		Update("restaurant_paid_at", now).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update orders"})
	}

	if err := tx.Save(&payout).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update payout status"})
	}

	// Mark all associated orders as paid
	payoutIDStr := strconv.Itoa(payout.PayoutID)
	if err := tx.Model(&models.Order{}).
		Where("restaurant_payout_id = ?", payoutIDStr).
		Update("restaurant_paid_at", now).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to mark orders as paid"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payout verified successfully"})
}

func GetRestaurantPayouts(c *fiber.Ctx) error {

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
	var restaurantPayouts []models.RestaurantPayout
	if err := database.DB.Preload("Shop").Limit(limit).Offset((page - 1) * limit).Order("created_at DESC").Find(&restaurantPayouts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch restaurant payouts"})
	}

	return c.Status(fiber.StatusOK).JSON(restaurantPayouts)
}

func FetchRestaurantPayouts(c *fiber.Ctx) error {

	user_id := c.Locals("user_id")
	shop_id := c.Params("shop_id")
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var restaurantPayouts []models.RestaurantPayout
	if err := database.DB.Preload("Shop").Where("shop_id = ?", shop_id).Limit(limit).Offset(offset).Order("created_at DESC").Find(&restaurantPayouts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch restaurant payouts"})
	}

	return c.Status(fiber.StatusOK).JSON(restaurantPayouts)
}
