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

	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}
	if err := c.Bind().JSON(&wallet_transaction).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if err := database.DB.First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	wallet_transaction.UserID = user_id
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
	if err := tx.First(&wallet, "user_id = ?", walletTransaction.UserID).Error; err != nil {
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
	if err := database.DB.Where("user_id = ?", user_id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"wallet": wallet})
}
