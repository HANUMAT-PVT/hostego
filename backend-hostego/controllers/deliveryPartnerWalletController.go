package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"
	"strconv"

	"github.com/gofiber/fiber/v3"
)

func FetchDeliveryPartnerWallet(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	var deliveryPartnerWallet models.DeliveryPartnerWallet

	if middleErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   middleErr.Error(),
		})
	}

	err := database.DB.Where("user_id = ?", user_id).First(&deliveryPartnerWallet).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(deliveryPartnerWallet)

}

func FetchDeliveryPartnerWalletTransactions(c fiber.Ctx) error {
	queryPageLimit := c.Query("limit", "10")
	queryPage := c.Query("page", "1")

	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   middleErr.Error(),
		})
	}
	pageLimit, err := strconv.Atoi(queryPageLimit)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid page limit",
		})
	}
	page, err := strconv.Atoi(queryPage)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid page",
		})
	}

	var deliveryPartnerWalletTransactions []models.DeliveryPartnerWalletTransaction
	err = database.DB.Where("user_id = ?", user_id).Order("created_at DESC").Limit(pageLimit).Offset(page).Find(&deliveryPartnerWalletTransactions).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Internal Server Error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(deliveryPartnerWalletTransactions)

}


func AddEarningsToDeliveryPartnerWallet(c fiber.Ctx, currentOrder models.Order) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   middleErr.Error(),
		})
	}
	if user_id == "" {
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

	var deliveryPartnerWallet models.DeliveryPartnerWallet
	var deliveryPartnerWalletTransaction models.DeliveryPartnerWalletTransaction

	err := tx.Where("user_id = ?", user_id).First(&deliveryPartnerWallet).Error
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Delivery Partner Wallet Not Found",
		})
	}

	deliveryPartnerWallet.Balance += currentOrder.DeliveryPartnerFee
	err = tx.Save(&deliveryPartnerWallet).Error
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Wallet Update Failed",
		})
	}
	deliveryPartnerWalletTransaction.DeliveryPartnerId = currentOrder.DeliveryPartnerId
	deliveryPartnerWalletTransaction.Amount = currentOrder.DeliveryPartnerFee
	deliveryPartnerWalletTransaction.TransactionType = models.TransactionCustomType(models.TransactionCredit)
	deliveryPartnerWalletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)

	err = tx.Create(&deliveryPartnerWalletTransaction).Error
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Wallet Transaction Creation Failed",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Delivery Partner Wallet Updated Successfully",
		"data":    deliveryPartnerWallet,
	})
}

func CreateWalletWithdrawalRequests(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
			"error":   middleErr.Error(),
		})
	}
	if user_id == "" {
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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch delivery partner wallets",
			"error":   err.Error(),
		})
	}

	if len(deliveryPartnerWallets) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "No delivery partners found with balance greater than 0",
		})
	}

	var createdTransactions []models.DeliveryPartnerWalletTransaction

	// Create withdrawal requests for each eligible delivery partner
	for _, wallet := range deliveryPartnerWallets {
		var deliveryPartner models.DeliveryPartner
		err := tx.Where("user_id = ?", wallet.UserId).First(&deliveryPartner).Error
		if err != nil {
			continue // Skip if delivery partner not found
		}

		transaction := models.DeliveryPartnerWalletTransaction{
			DeliveryPartnerId: deliveryPartner.DeliveryPartnerID.String(),
			Amount:            wallet.Balance,
			TransactionType:   models.TransactionCustomType(models.TransactionDebit),
			TransactionStatus: models.TransactionStatusType(models.TransactionPending),
			UserId:            wallet.UserId,

		}

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to create wallet transactions",
				"error":   err.Error(),
			})
		}

		createdTransactions = append(createdTransactions, transaction)
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to commit transaction",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Wallet Withdrawal Requests Created Successfully",
		"data": fiber.Map{
			"transactions_created": len(createdTransactions),
			"transactions":         createdTransactions,
		},
	})
}

func VerifyDeliveryPartnerWithdrawalRequest(c fiber.Ctx) error {
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
	var requestData struct {
		UniqueTransactionID string `json:"unique_transaction_id"`
	}
	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
	}

	transactionId := c.Params("transaction_id")

	var deliveryPartnerWalletTransaction models.DeliveryPartnerWalletTransaction
	deliveryPartnerWalletTransaction.PaymentMethod.PaymentVerifiedByAdmin = user_id
	deliveryPartnerWalletTransaction.PaymentMethod.UniqueTransactionID = requestData.UniqueTransactionID
	deliveryPartnerWalletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)
	err := tx.Where("transaction_id = ?", transactionId).Save(&deliveryPartnerWalletTransaction).Error

	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to update transaction status",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Delivery Partner Wallet Transaction Status Updated Successfully",
		"data":    deliveryPartnerWalletTransaction,
	})
}
