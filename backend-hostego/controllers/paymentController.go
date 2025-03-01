package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"
	"time"

	"github.com/gofiber/fiber/v3"
)

func InitiatePayment(c fiber.Ctx) error {
	userId, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}

	type OrderRequest struct {
		OrderID string `json:"order_id"` //Only accept Order ID, not amount
	}
	var order models.Order

	var request OrderRequest

	var cartItem models.CartItem

	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	var wallet models.Wallet

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.First(&order, "order_id=?", request.OrderID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if err := tx.First(&wallet, "user_id=?", userId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if wallet.Balance < order.FinalOrderValue {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Wallet balance insufficent to complete payment"})
	}
	totalAmountToDeduct := order.FinalOrderValue
	wallet.Balance -= totalAmountToDeduct
	var walletTransaction models.WalletTransaction

	walletTransaction.Amount = totalAmountToDeduct
	walletTransaction.TransactionType = "debit"
	walletTransaction.UserId = userId
	walletTransaction.TransactionStatus = "success"

	var paymentTransaction models.PaymentTransaction

	if err := tx.Create(&walletTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	paymentTransaction.OrderId = request.OrderID
	paymentTransaction.UserId = userId
	paymentTransaction.Amount = totalAmountToDeduct
	paymentTransaction.PaymentStatus = "success"
	paymentTransaction.PaymentMethod = "wallet"
	order.PaymentTransactionId = paymentTransaction.PaymentTransactionId
	order.OrderStatus = "placed"

	if err := tx.Create(&paymentTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	order.PaymentTransactionId = paymentTransaction.PaymentTransactionId
	if err := tx.Where("user_id=?", userId).Save(&wallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	if err := tx.Preload("User").Where("order_id=?", request.OrderID).Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if err := tx.Where("user_id = ?", userId).Delete(&cartItem).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove cart items"})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment Completed", "payment_transaction": paymentTransaction, "order": order, "wallet_transaction": walletTransaction})
}

func FetchUserPaymentTransactions(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}

	var payment_transactions []models.PaymentTransaction

	if err := database.DB.Find(&payment_transactions, "user_id=?", user_id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(payment_transactions)
}

func InitiateRefundPayment(c fiber.Ctx) error {
	current_user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	type OrderRequest struct {
		OrderID string `json:"order_id"`
		UserId  string `json:"user_id"`
	}

	var request OrderRequest

	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if current_user_id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	var order models.Order
	var wallet models.Wallet
	var walletTransaction models.WalletTransaction

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := database.DB.Where("order_id = ?", request.OrderID).First(&order).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := tx.Where("user_id=?", request.UserId).First(&wallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	order.OrderStatus = models.OrderStatusType(models.CanceledOrderStatus)
	order.Refunded = true
	order.RefundedAt = time.Now()
	order.RefundInitiator = current_user_id
	order.DeliveryPartnerId = ""
	order.DeliveryPartner = nil

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := tx.Where("user_id=?", request.UserId).First(&wallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	walletTransaction.Amount = order.FinalOrderValue
	walletTransaction.TransactionType = models.TransactionCustomType(models.TransactionRefund)
	walletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)
	walletTransaction.UserId = request.UserId
	wallet.Balance += order.FinalOrderValue

	if err := tx.Create(&walletTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := tx.Where("user_id=?", request.UserId).Save(&wallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Refund Completed", "wallet_transaction": walletTransaction, "wallet": wallet})
}
