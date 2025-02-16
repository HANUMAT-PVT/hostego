   package controllers

   import (
   	"backend-hostego/database"
   	"backend-hostego/middlewares"
   	"backend-hostego/models"

   	"github.com/gofiber/fiber/v3"
   )

   func InitiatePayment(c fiber.Ctx) error {
   	userId, middleErr := middlewares.VerifyUserAuthCookie(c)
   	if middleErr != nil {
   		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
   	}

   	type OrderRequest struct {
   		OrderID string `json:"order_id"`    //Only accept Order ID, not amount
   	}
   	var order models.Order

   	var request OrderRequest

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
   		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
   	}

	totalAmountToDeduct:=order.FinalOrderValue;

	

   	var paymentTransaction models.PaymentTransaction

   	if err := tx.First(&paymentTransaction, "transaction_id = ?", walletTransactionID).Error; err != nil {
   		tx.Rollback()
   		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Transaction not found"})
   	}

   	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment Completed", "payment_transaction": paymentTransaction})
   }
