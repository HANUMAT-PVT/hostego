package controllers

import (
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"

	"github.com/gofiber/fiber/v3"
)

func CreditWalletTransaction(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	var user models.User
	var wallet_transaction models.WalletTransaction
	var requestData struct {
		Amount                  float64 `json:"amount"`
		PaymentScreenShotImgUrl string  `json:"payment_screenshot_img_url"`
		UniqueTransactionID     string  `json:"unique_transaction_id"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	// if err := c.Bind().JSON(&wallet_transaction).Error; err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	// }
	if err := database.DB.First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	wallet_transaction.UserId = user_id
	wallet_transaction.TransactionType = "credit"
	wallet_transaction.TransactionStatus = models.TransactionPending
	wallet_transaction.Amount = requestData.Amount
	wallet_transaction.PaymentMethod.PaymentScreenShotImgUrl = requestData.PaymentScreenShotImgUrl
	wallet_transaction.PaymentMethod.UniqueTransactionID = requestData.UniqueTransactionID

	if err := database.DB.Create(&wallet_transaction).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Wallet Transaction Created"})
}

func VerifyWalletTransactionById(c fiber.Ctx) error {
	userID, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}

	walletTransactionID := c.Params("id")

	// Start a transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var walletTransaction models.WalletTransaction
	if err := tx.First(&walletTransaction, "transaction_id = ?", walletTransactionID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Transaction not found"})
	}

	var wallet models.Wallet
	if err := tx.First(&wallet, "user_id = ?", walletTransaction.UserId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Wallet not found"})
	}

	// Update wallet balance
	wallet.Balance += walletTransaction.Amount
	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update wallet balance"})
	}

	// Update transaction status
	walletTransaction.TransactionStatus = models.TransactionSuccess
	walletTransaction.PaymentMethod.PaymentVerifiedByAdmin = userID
	if err := tx.Save(&walletTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update transaction status"})
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Wallet transaction verified successfully"})
}

func FetchUserWallet(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	var wallet models.Wallet
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if err := database.DB.Preload("User").Where("user_id = ?", user_id).Find(&wallet).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(wallet)
}

func FetchUserWalletTransactions(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}

	var wallet_transactions []models.WalletTransaction

	if err := database.DB.Find(&wallet_transactions, "user_id=?", user_id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(wallet_transactions)
}

func FetchAllWalletTransactions(c fiber.Ctx) error {
	user_id, middleErr := middlewares.VerifyUserAuthCookie(c)

	dbQuery := database.DB

	status := c.Query("status")
	// tagsQuery := c.Query("tags")

	if status != "" {
		dbQuery = dbQuery.Where("transaction_status = ?", status)
	}
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if user_id != "admin" {
		// return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	var wallet_transactions []models.WalletTransaction
	if err := dbQuery.Find(&wallet_transactions).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(wallet_transactions)
}
