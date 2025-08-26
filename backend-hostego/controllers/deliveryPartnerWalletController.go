package controllers

import (
	"backend-hostego/database"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FetchDeliveryPartnerWallet(c *fiber.Ctx) error {
	deliveryPartnerId := c.Params("id")
	user_id := c.Locals("user_id")
	var deliveryPartnerWallet models.DeliveryPartnerWallet

	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	err := database.DB.Where("delivery_partner_id = ?", deliveryPartnerId).First(&deliveryPartnerWallet).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(deliveryPartnerWallet)

}

func FetchDeliveryPartnerWalletTransactions(c *fiber.Ctx) error {

	deliveryPartnerId := c.Params("id")
	user_id := c.Locals("user_id")

	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var deliveryPartnerWalletTransactions []models.DeliveryPartnerWalletTransaction

	err := database.DB.Where("delivery_partner_id = ?", deliveryPartnerId).Order("created_at DESC").Find(&deliveryPartnerWalletTransactions).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(deliveryPartnerWalletTransactions)

}

func AddEarningsToDeliveryPartnerWallet(currentOrder models.Order) error {
	return database.SafeTransactionWithCleanup(func(tx *gorm.DB) error {
		var deliveryPartnerWallet models.DeliveryPartnerWallet
		var deliveryPartnerWalletTransaction models.DeliveryPartnerWalletTransaction

		// Fetch delivery partner wallet
		if err := tx.Where("delivery_partner_id = ?", currentOrder.DeliveryPartnerId).First(&deliveryPartnerWallet).Error; err != nil {
			return err
		}

		// Update wallet balance
		deliveryPartnerWallet.Balance += currentOrder.DeliveryPartnerFee
		if err := tx.Save(&deliveryPartnerWallet).Error; err != nil {
			return err
		}

		// Create wallet transaction record
		deliveryPartnerWalletTransaction.Amount = currentOrder.DeliveryPartnerFee
		deliveryPartnerWalletTransaction.TransactionType = models.TransactionCustomType(models.TransactionCredit)
		deliveryPartnerWalletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)
		deliveryPartnerWalletTransaction.DeliveryPartnerId = currentOrder.DeliveryPartnerId

		if err := tx.Create(&deliveryPartnerWalletTransaction).Error; err != nil {
			return err
		}

		return nil
	}, "AddEarningsToDeliveryPartnerWallet")
}

func CreateWalletWithdrawalRequests(c *fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var createdTransactions []models.DeliveryPartnerWalletTransaction

	err := database.SafeTransactionWithCleanup(func(tx *gorm.DB) error {
		// Find all delivery partners with wallet balance > 0
		var deliveryPartnerWallets []models.DeliveryPartnerWallet
		if err := tx.Where("balance > ?", 0).Find(&deliveryPartnerWallets).Error; err != nil {
			return err
		}

		if len(deliveryPartnerWallets) == 0 {
			return nil // No wallets with balance to process
		}

		// Create withdrawal requests for each eligible delivery partner
		for _, wallet := range deliveryPartnerWallets {
			var deliveryPartner models.DeliveryPartner
			if err := tx.Where("delivery_partner_id = ?", wallet.DeliveryPartnerId).First(&deliveryPartner).Error; err != nil {
				continue // Skip if delivery partner not found
			}

			transaction := models.DeliveryPartnerWalletTransaction{
				Amount:            wallet.Balance,
				TransactionType:   models.TransactionCustomType(models.TransactionDebit),
				TransactionStatus: models.TransactionStatusType(models.TransactionPending),
				DeliveryPartnerId: wallet.DeliveryPartnerId,
				DeliveryPartner:   deliveryPartner,
			}

			if err := tx.Create(&transaction).Error; err != nil {
				return err
			}

			createdTransactions = append(createdTransactions, transaction)
		}

		return nil
	}, "CreateWalletWithdrawalRequests")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to create withdrawal requests",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Withdrawal requests created successfully",
		"data":    createdTransactions,
	})
}

func VerifyDeliveryPartnerWithdrawalRequest(c *fiber.Ctx) error {
	transactionId := c.Params("transaction_id")
	user_id := c.Locals("user_id").(int)
	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var requestData struct {
		UniqueTransactionID          string `json:"unique_transaction_id"`
		TransactionStatusTypePayment string `json:"transaction_status"`
	}
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	var result struct {
		Transaction models.DeliveryPartnerWalletTransaction
		Wallet      models.DeliveryPartnerWallet
	}

	err := database.SafeTransactionWithCleanup(func(tx *gorm.DB) error {
		// First load the existing transaction
		if err := tx.Where("transaction_id = ?", transactionId).First(&result.Transaction).Error; err != nil {
			return err
		}

		// Load the wallet
		if err := tx.Where("delivery_partner_id = ?", result.Transaction.DeliveryPartnerId).First(&result.Wallet).Error; err != nil {
			return err
		}

		// Update transaction
		result.Transaction.PaymentMethod.PaymentVerifiedByAdmin = user_id
		result.Transaction.TransactionStatus = models.TransactionStatusType(requestData.TransactionStatusTypePayment)

		// Update wallet balance if transaction is successful
		if models.TransactionStatusType(requestData.TransactionStatusTypePayment) == models.TransactionSuccess {
			result.Wallet.Balance -= result.Transaction.Amount
			result.Transaction.PaymentMethod.UniqueTransactionID = requestData.UniqueTransactionID
		}

		// Save both updates
		if err := tx.Save(&result.Transaction).Error; err != nil {
			return err
		}

		if err := tx.Save(&result.Wallet).Error; err != nil {
			return err
		}

		return nil
	}, "VerifyDeliveryPartnerWithdrawalRequest")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to verify withdrawal request",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Withdrawal request verified successfully",
		"data":    result.Transaction,
	})
}

func FetchAllDeliveryPartnersTransactions(c *fiber.Ctx) error {

	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var deliveryParnterWalletTransactions []models.DeliveryPartnerWalletTransaction
	err := database.DB.Where("transaction_type = ?", "debit").Order("created_at asc").Preload("DeliveryPartner.User").Find(&deliveryParnterWalletTransactions).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(deliveryParnterWalletTransactions)

}
