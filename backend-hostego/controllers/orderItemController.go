package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CancelOrderItemAndInitiateRefund(c *fiber.Ctx) error {
	// Safe user ID extraction
	current_user_id, err := middlewares.SafeUserIDExtractor(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user authentication: " + err.Error()})
	}
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

	var result struct {
		Order             models.Order
		Wallet            models.Wallet
		WalletTransaction models.WalletTransaction
		Product           models.Product
		OrderItem         models.OrderItem
	}

	err = database.SafeTransactionWithCleanup(func(tx *gorm.DB) error {
		// Fetch order
		if err := tx.Where("order_id = ?", request.OrderID).First(&result.Order).Error; err != nil {
			return err
		}

		// Fetch wallet
		if err := tx.Where("user_id=?", result.Order.UserId).First(&result.Wallet).Error; err != nil {
			return err
		}

		// Fetch order item
		if err := tx.Where("product_id=? AND order_id=? AND user_id=?", request.ProductId, request.OrderID, result.Order.UserId).First(&result.OrderItem).Error; err != nil {
			return err
		}

		// Fetch product
		if err := tx.Where("product_id=?", request.ProductId).First(&result.Product).Error; err != nil {
			return err
		}

		// Create wallet transaction for refund
		result.WalletTransaction.TransactionType = models.TransactionCustomType(models.TransactionRefund)
		result.WalletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)
		result.WalletTransaction.UserId = result.Order.UserId

		if request.Quantity < result.OrderItem.Quantity {
			if request.Quantity <= 0 {
				// Full refund - delete order item
				result.WalletTransaction.Amount = result.OrderItem.SubTotal
				result.Wallet.Balance += result.OrderItem.SubTotal

				if err := tx.Where("order_item_id=?", result.OrderItem.OrderItemId).Delete(&result.OrderItem).Error; err != nil {
					return err
				}
			} else {
				// Partial refund - update order item
				var amountToRefund = result.Product.FoodPrice * float64(request.Quantity)
				result.WalletTransaction.Amount = amountToRefund
				result.Wallet.Balance += amountToRefund
				result.OrderItem.Quantity = request.Quantity
				result.OrderItem.SubTotal = result.OrderItem.SubTotal - amountToRefund

				if err := tx.Where("order_item_id=?", result.OrderItem.OrderItemId).Save(&result.OrderItem).Error; err != nil {
					return err
				}
			}
		}

		// Create wallet transaction
		if err := tx.Create(&result.WalletTransaction).Error; err != nil {
			return err
		}

		// Update wallet balance
		if err := tx.Where("user_id=?", result.Order.UserId).Save(&result.Wallet).Error; err != nil {
			return err
		}

		return nil
	}, "CancelOrderItemAndInitiateRefund")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":            "Refund Completed",
		"wallet_transaction": result.WalletTransaction,
		"wallet":             result.Wallet,
	})
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
