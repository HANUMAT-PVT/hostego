package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v2"
)

func CancelOrderItemAndInitiateRefund(c *fiber.Ctx) error {

	current_user_id := c.Locals("user_id").(int)
	if current_user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	type OrderRequest struct {
		OrderID   int `json:"order_id"`
		Quantity  int `json:"quantity"`
		ProductId int `json:"product_id"`
	}

	var request OrderRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if current_user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	var order models.Order
	var wallet models.Wallet
	var walletTransaction models.WalletTransaction
	var product models.Product
	var orderItem models.OrderItem

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := database.DB.Where("order_id = ?", request.OrderID).First(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := tx.Where("user_id=?", order.UserId).First(&wallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := tx.Where("product_id=? AND order_id=? AND user_id=?", request.ProductId, request.OrderID, order.UserId).First(&orderItem).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := tx.Where("product_id=?", request.ProductId).First(&product).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// fetching user wallet

	if err := tx.Where("user_id=?", order.UserId).First(&wallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Create a wallet transaction for the refund
	walletTransaction.TransactionType = models.TransactionCustomType(models.TransactionRefund)
	walletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)
	walletTransaction.UserId = order.UserId

	if request.Quantity < orderItem.Quantity {
		if request.Quantity <= 0 {
			walletTransaction.Amount = orderItem.SubTotal
			wallet.Balance += orderItem.SubTotal

			if err := tx.Where("order_item_id=?", orderItem.OrderItemId).Delete(&orderItem).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
		} else {
			var amountToRefund = product.FoodPrice * float64(request.Quantity)
			walletTransaction.Amount = amountToRefund
			wallet.Balance += amountToRefund
			orderItem.Quantity = request.Quantity
			orderItem.SubTotal = orderItem.SubTotal - amountToRefund

			if err := tx.Where("order_item_id=?", orderItem.OrderItemId).Save(&orderItem).Error; err != nil {
				tx.Rollback()
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
		}

	}

	// create a wallet transaction for the refund
	if err := tx.Create(&walletTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// Update the wallet balance
	if err := tx.Where("user_id=?", order.UserId).Save(&wallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Refund Completed", "wallet_transaction": walletTransaction, "wallet": wallet})
}

func FetchOrderItems(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var query = database.DB
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate != "" && endDate != "" {
		query = query.Where("orders.created_at BETWEEN ? AND ?", startDate+" 00:00:00", endDate+" 23:59:59")
	}
	var orders []models.Order

	query.Preload("Product").Preload("User").Preload("Product.Shop").Where("order_status=?", "delivered").Order("created_at ASC").Find(&orders)

	return c.Status(fiber.StatusOK).JSON(orders)
}
