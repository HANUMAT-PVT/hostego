package controllers

import (
	"backend-hostego/config"
	"backend-hostego/database"
	"backend-hostego/middlewares"
	"backend-hostego/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	cashfree "github.com/cashfree/cashfree-pg/v3"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func InitiatePayment(c fiber.Ctx) error {
	userId := c.Locals("user_id").(int)
	if userId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type OrderRequest struct {
		OrderID int `json:"order_id"` //Only accept Order ID, not amount
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

	// Create order items from cart items
	var orderItems []models.CartItem
	if err := json.Unmarshal(order.OrderItems, &orderItems); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse order items"})
	}

	// Store each cart item as an order item
	for _, item := range orderItems {
		orderItem := models.OrderItem{
			OrderId:   order.OrderId,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			SubTotal:  item.SubTotal,
			UserId:    order.UserId,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}

		// Update product stock
		if err := tx.Model(&models.Product{}).
			Where("product_id = ?", item.ProductId).
			Update("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}
	}

	// Mark cart items as deleted
	if err := tx.
		Where("user_id = ?", userId).
		Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete cart items"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment Completed", "payment_transaction": paymentTransaction, "order": order, "wallet_transaction": walletTransaction})
}

func FetchUserPaymentTransactions(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var payment_transactions []models.PaymentTransaction

	if err := database.DB.Find(&payment_transactions, "user_id=?", user_id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(payment_transactions)
}

func InitiateRefundPayment(c fiber.Ctx) error {
	current_user_id := c.Locals("user_id").(int)
	if current_user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	type OrderRequest struct {
		OrderID int `json:"order_id"`
	}

	var request OrderRequest

	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}
	if current_user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	var order models.Order
	var wallet models.Wallet
	var walletTransaction models.WalletTransaction
	var delivery_partner models.DeliveryPartner

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
	if order.DeliveryPartnerId != 0 {
		if err := tx.Where("delivery_partner_id = ?", order.DeliveryPartnerId).First(&delivery_partner).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	order.OrderStatus = models.OrderStatusType(models.CanceledOrderStatus)
	order.Refunded = true
	order.RefundedAt = time.Now()
	order.RefundInitiator = current_user_id
	order.DeliveryPartnerId = 0
	order.DeliveryPartner = nil

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := tx.Where("user_id=?", order.UserId).First(&wallet).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Create a wallet transaction for the refund
	walletTransaction.Amount = order.FinalOrderValue
	walletTransaction.TransactionType = models.TransactionCustomType(models.TransactionRefund)
	walletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)
	walletTransaction.UserId = order.UserId
	wallet.Balance += order.FinalOrderValue

	var orderItems []models.CartItem
	if err := json.Unmarshal(order.OrderItems, &orderItems); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse order items"})
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
	for _, item := range orderItems {
		if err := tx.Model(&models.Product{}).Where("product_id = ?", item.ProductId).
			Update("stock_quantity", gorm.Expr("stock_quantity + ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := tx.Where("order_id = ?", order.OrderId).
			Delete(&models.OrderItem{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}
	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Refund Completed", "wallet_transaction": walletTransaction, "wallet": wallet})
}

func InitateCashfreePaymentOrder(c fiber.Ctx) error {

	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var orderRequest struct {
		OrderId int `json:"order_id"`
	}
	err := c.Bind().JSON(&orderRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println(orderRequest, "orderRequest", user_id)
	var user models.User
	if err := tx.First(&user, "user_id = ?", user_id).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	var order models.Order
	if err := tx.First(&order, "order_id = ?", orderRequest.OrderId).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	fmt.Println("Cashfree Client ID:", config.GetEnv("CASHFREE_CLIENT_ID_"))
	fmt.Println("Cashfree Client Secret:", config.GetEnv("CASHFREE_CLIENT_SECRET_"))

	clientId := config.GetEnv("CASHFREE_CLIENT_ID_")
	clientSecret := config.GetEnv("CASHFREE_CLIENT_SECRET_")
	cashfree.XClientId = &clientId
	cashfree.XClientSecret = &clientSecret
	cashfree.XEnvironment = cashfree.SANDBOX

	request := cashfree.CreateOrderRequest{
		OrderAmount: order.FinalOrderValue,
		CustomerDetails: cashfree.CustomerDetails{
			CustomerId:    strconv.Itoa(user.UserId),
			CustomerPhone: user.MobileNumber,
			CustomerEmail: &user.Email,
			CustomerName:  &user.FirstName,
		},
		OrderCurrency: "INR",
		OrderSplits:   []cashfree.VendorSplit{},
	}
	version := "2023-08-01"
	response, httpResponse, err := cashfree.PGCreateOrder(&version, &request, nil, nil, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(httpResponse.StatusCode)
	fmt.Println(response)
	return c.Status(httpResponse.StatusCode).JSON(response)

}

func VerifyCashfreePayment(c fiber.Ctx) error {
	user_id := c.Locals("user_id")
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	cf_order_id := c.Params("cf_order_id")

	version := "2023-08-01"
	response, httpResponse, err := cashfree.PGFetchOrder(&version, cf_order_id, nil, nil, nil)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(httpResponse.StatusCode).JSON(fiber.Map{"erro": err.Error()})
	}
	if *response.OrderStatus == "ACTIVE" {
		fmt.Println("ACTIVE")
	}
	if *response.OrderStatus == "PAID" {
		fmt.Println("PAID")
	}
	userId, middleErr := middlewares.VerifyUserAuthCookie(c)
	if middleErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": middleErr.Error()})
	}

	type OrderRequest struct {
		OrderID int `json:"order_id"` //Only accept Order ID, not amount
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

	// Create order items from cart items
	var orderItems []models.CartItem
	if err := json.Unmarshal(order.OrderItems, &orderItems); err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse order items"})
	}

	// Store each cart item as an order item
	for _, item := range orderItems {
		orderItem := models.OrderItem{
			OrderId:   order.OrderId,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
			SubTotal:  item.SubTotal,
			UserId:    order.UserId,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}

		// Update product stock
		if err := tx.Model(&models.Product{}).
			Where("product_id = ?", item.ProductId).
			Update("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}
	}

	// Mark cart items as deleted
	if err := tx.
		Where("user_id = ?", userId).
		Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete cart items"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment Completed", "payment_transaction": paymentTransaction, "order": order, "wallet_transaction": walletTransaction})

}
