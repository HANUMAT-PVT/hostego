package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func FetchDeliveryPartnerWallet(c fiber.Ctx) error {
	deliveryPartnerId := c.Params("id")
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	var deliveryPartnerWallet models.DeliveryPartnerWallet

	if middleErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   middleErr.Error(),
		})
	}
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

func FetchDeliveryPartnerWalletTransactions(c fiber.Ctx) error {

	deliveryPartnerId := c.Params("id")
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   middleErr.Error(),
		})
	}
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

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var deliveryPartnerWallet models.DeliveryPartnerWallet
	var deliveryPartnerWalletTransaction models.DeliveryPartnerWalletTransaction

	err := tx.Where("delivery_partner_id = ?", currentOrder.DeliveryPartnerId).First(&deliveryPartnerWallet).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	deliveryPartnerWallet.Balance += currentOrder.DeliveryPartnerFee
	err = tx.Save(&deliveryPartnerWallet).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	deliveryPartnerWalletTransaction.Amount = currentOrder.DeliveryPartnerFee
	deliveryPartnerWalletTransaction.TransactionType = models.TransactionCustomType(models.TransactionCredit)
	deliveryPartnerWalletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)
	deliveryPartnerWalletTransaction.DeliveryPartnerId = currentOrder.DeliveryPartnerId
	err = tx.Create(&deliveryPartnerWalletTransaction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func CreateWalletWithdrawalRequests(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return middleErr
	}
	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Find all delivery partners with wallet balance > 0
	var deliveryPartnerWallets []models.DeliveryPartnerWallet
	err := tx.Where("balance > ?", 0).Find(&deliveryPartnerWallets).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	if len(deliveryPartnerWallets) == 0 {
		return err
	}

	var createdTransactions []models.DeliveryPartnerWalletTransaction

	// Create withdrawal requests for each eligible delivery partner
	for _, wallet := range deliveryPartnerWallets {
		var deliveryPartner models.DeliveryPartner
		err := tx.Where("delivery_partner_id = ?", wallet.DeliveryPartnerId).First(&deliveryPartner).Error
		if err != nil {
			continue // Skip if delivery partner not found
		}

		transaction := models.DeliveryPartnerWalletTransaction{

			Amount:            wallet.Balance,
			TransactionType:   models.TransactionCustomType(models.TransactionDebit),
			TransactionStatus: models.TransactionStatusType(models.TransactionPending),
			DeliveryPartnerId: wallet.DeliveryPartnerId,
		}

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			return err
		}

		createdTransactions = append(createdTransactions, transaction)
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Withdrawal requests created successfully",
		"data":    createdTransactions,
	})
}

func VerifyDeliveryPartnerWithdrawalRequest(c fiber.Ctx) error {
	transactionId := c.Params("transaction_id")
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   middleErr.Error(),
		})
	}

	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// First load the existing transaction
	var deliveryPartnerWalletTransaction models.DeliveryPartnerWalletTransaction
	if err := tx.Where("transaction_id = ?", transactionId).First(&deliveryPartnerWalletTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Transaction not found",
			"error":   err.Error(),
		})
	}

	var deliveryPartnerWallet models.DeliveryPartnerWallet
	if err := tx.Where("delivery_partner_id = ?", deliveryPartnerWalletTransaction.DeliveryPartnerId).First(&deliveryPartnerWallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Wallet not found",
			"error":   err.Error(),
		})
	}

	var requestData struct {
		UniqueTransactionID string `json:"unique_transaction_id"`
	}
	if err := c.Bind().JSON(&requestData); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	// Update transaction
	deliveryPartnerWalletTransaction.PaymentMethod.PaymentVerifiedByAdmin = user_id
	deliveryPartnerWalletTransaction.PaymentMethod.UniqueTransactionID = requestData.UniqueTransactionID
	deliveryPartnerWalletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)

	// Update wallet balance
	deliveryPartnerWallet.Balance -= deliveryPartnerWalletTransaction.Amount

	// Save both updates
	if err := tx.Save(&deliveryPartnerWalletTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update transaction",
			"error":   err.Error(),
		})
	}

	if err := tx.Save(&deliveryPartnerWallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update wallet",
			"error":   err.Error(),
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to commit transaction",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Withdrawal request verified successfully",
		"data":    deliveryPartnerWalletTransaction,
	})
}

func FetchAllDeliveryPartnersTransactions(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Found the delivery Partners",
			"err":     middleErr.Error(),
		})
	}
	if user_id == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	var deliveryParnterWalletTransactions []models.DeliveryPartnerWalletTransaction
	err := database.DB.Find(&deliveryParnterWalletTransactions).Error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Something went wrong",
			"err":     err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(deliveryParnterWalletTransactions)

}
