package controllers

import (
	"backend-hostego/config"
	"backend-hostego/database"
	"backend-hostego/models"
	natsclient "backend-hostego/nats"
	"crypto/hmac"
	"crypto/sha256"
	"strconv"
	"time"

	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	razorpay "github.com/razorpay/razorpay-go"
	"gorm.io/gorm"
)

var rz_key_id = config.GetEnv("RAZORPAY_KEY_ID_")
var rz_key_secret = config.GetEnv("RAZORPAY_KEY_SECRET_")
var rz_client = razorpay.NewClient(rz_key_id, rz_key_secret)

func InitiatePayment(c *fiber.Ctx) error {
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

	if err := c.BodyParser(&request); err != nil {
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
	natsclient.SendMessageToUsersByRole(orderManagerRoles, "New Order Placed", "Please check the details and take the necessary action.")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment Completed", "payment_transaction": paymentTransaction, "order": order, "wallet_transaction": walletTransaction})
}

func FetchUserPaymentTransactions(c *fiber.Ctx) error {
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

func InitiateRefundPayment(c *fiber.Ctx) error {
	current_user_id := c.Locals("user_id").(int)
	if current_user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}
	type OrderRequest struct {
		OrderID int `json:"order_id"`
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

func InitateCashfreePaymentOrder(c *fiber.Ctx) error {

	user_id := c.Locals("user_id").(int)
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
	var paymentTransaction models.PaymentTransaction

	err := c.BodyParser(&orderRequest)
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

	body := map[string]interface{}{
		"order_amount":   order.FinalOrderValue,
		"order_currency": "INR",
		"customer_details": map[string]interface{}{
			"customer_id":    strconv.Itoa(user.UserId),
			"customer_phone": user.MobileNumber,
			"customer_email": user.Email,
			"customer_name":  user.FirstName,
		},
	}

	restyClient := resty.New()
	// Cashfree credentials from env
	clientId := config.GetEnv("CASHFREE_CLIENT_ID_")
	clientSecret := config.GetEnv("CASHFREE_CLIENT_SECRET_")
	cashFreeApiUrl := config.GetEnv("CASHFREE_API_URL_")

	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-version", "2023-08-01").
		SetHeader("x-client-id", clientId).
		SetHeader("x-client-secret", clientSecret).
		SetBody(body).
		Post(cashFreeApiUrl)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Return response from Cashfree
	var cashfreeResp map[string]interface{}
	if err := json.Unmarshal(resp.Body(), &cashfreeResp); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid response from Cashfree"})
	}
	paymentTransaction.OrderId = orderRequest.OrderId
	paymentTransaction.UserId = user_id
	paymentTransaction.Amount = order.FinalOrderValue
	paymentTransaction.PaymentStatus = "pending"
	paymentTransaction.PaymentMethod = "UPI"
	paymentTransaction.PaymentOrderId = cashfreeResp["order_id"].(string)

	if err := tx.Create(&paymentTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	order.PaymentTransactionId = paymentTransaction.PaymentTransactionId

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(cashfreeResp)

}

func VerifyCashfreePayment(c *fiber.Ctx) error {

	user_id := c.Locals("user_id").(int)
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type OrderRequest struct {
		OrderID int `json:"order_id"` //Only accept Order ID, not amount
	}

	var order models.Order

	var request OrderRequest

	var cartItem models.CartItem

	// cf_order_id := c.Params("cf_order_id")

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var paymentTransaction models.PaymentTransaction

	if err := tx.Where("order_id=?", request.OrderID).First(&paymentTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	restyClient := resty.New()

	// Cashfree credentials from env

	clientId := config.GetEnv("CASHFREE_CLIENT_ID_")
	clientSecret := config.GetEnv("CASHFREE_CLIENT_SECRET_")
	cashFreeApiUrl := config.GetEnv("CASHFREE_API_URL_")

	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("x-api-version", "2023-08-01").
		SetHeader("x-client-id", clientId).
		SetHeader("x-client-secret", clientSecret).
		Post(cashFreeApiUrl + "/" + paymentTransaction.PaymentOrderId)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	// Return response from Cashfree
	var cashfreeResp map[string]interface {
	}
	if err := json.Unmarshal(resp.Body(), &cashfreeResp); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid response from Cashfree"})
	}
	// return c.Status(fiber.StatusOK).JSON(cashfreeResp)
	if cashfreeResp["order_status"] != "PAID" {
		return c.Status(500).JSON(fiber.Map{"error": "`Payment is not paid yet", "response": cashfreeResp})
	}

	if err := tx.First(&order, "order_id=?", request.OrderID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if order.OrderStatus != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Order is already Verifed and Placed !"})
	}

	totalAmountToDeduct := order.FinalOrderValue

	var walletTransaction models.WalletTransaction

	walletTransaction.Amount = totalAmountToDeduct
	walletTransaction.TransactionType = "debit"
	walletTransaction.UserId = user_id
	walletTransaction.TransactionStatus = "success"

	if err := tx.Create(&walletTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	paymentTransaction.PaymentStatus = "success"
	order.OrderStatus = "placed"

	if err := tx.Save(&paymentTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	order.PaymentTransactionId = paymentTransaction.PaymentTransactionId

	if err := tx.Preload("User").Where("order_id=?", request.OrderID).Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if err := tx.Where("user_id = ?", user_id).Delete(&cartItem).Error; err != nil {
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
		Where("user_id = ?", user_id).
		Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete cart items"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}
	natsclient.SendMessageToUsersByRole(orderManagerRoles, "New Order Placed", "Please check the details and take the necessary action.")
	log.Print("Payload sent to the frontend")

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment Completed", "payment_transaction": paymentTransaction, "order": order, "wallet_transaction": walletTransaction, "response": cashfreeResp})

}

func InitateRazorpayPaymentOrder(c *fiber.Ctx) error {

	user_id := c.Locals("user_id").(int)
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
	var paymentTransaction models.PaymentTransaction

	err := c.BodyParser(&orderRequest)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

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

	data := map[string]interface{}{
		"amount":   order.FinalOrderValue * 100, // Amount is in currency subunits. Default currency is INR. Hence, 50000 refers to 50000 paise
		"currency": "INR",
		"receipt":  "some_receipt_id",
	}
	fmt.Print("rz_client", rz_key_id, "client secret", rz_key_secret)
	body, err := rz_client.Order.Create(data, nil)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create order", "message": err.Error()})
	}
	fmt.Print("created order")
	paymentTransaction.OrderId = orderRequest.OrderId
	paymentTransaction.UserId = user_id
	paymentTransaction.Amount = order.FinalOrderValue
	paymentTransaction.PaymentStatus = "pending"
	paymentTransaction.PaymentMethod = "UPI"
	paymentTransaction.PaymentOrderId = body["id"].(string)

	if err := tx.Create(&paymentTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	order.PaymentTransactionId = paymentTransaction.PaymentTransactionId

	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"order_id": body["id"],
		"amount":   order.FinalOrderValue,
		"currency": "INR",
		"key":      rz_key_id,
	})

}

func VerifyRazorpayPayment(c *fiber.Ctx) error {

	type OrderRequest struct {
		OrderID           int    `json:"order_id"` //Only accept Order ID, not amount
		RazorpayOrderID   string `json:"razorpay_order_id"`
		PaymentID         string `json:"razorpay_payment_id"`
		RazorpaySignature string `json:"razorpay_signature"`
	}

	var order models.Order

	var request OrderRequest

	var cartItem models.CartItem

	// cf_order_id := c.Params("cf_order_id")

	tx := database.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var paymentTransaction models.PaymentTransaction

	if err := tx.Where("order_id=?", request.OrderID).First(&paymentTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	// here razorpay secret

	if VerifyRazorpaySignature(request.RazorpayOrderID, request.PaymentID, request.RazorpaySignature, rz_key_secret) != true {

		return c.Status(400).JSON(fiber.Map{"error": "Signature verification failed"})
	}

	if err := tx.First(&order, "order_id=?", request.OrderID).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if order.OrderStatus != "pending" {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order is already Verifed and Placed !", "response": fiber.Map{"order_status": "PAID"}})
	}

	totalAmountToDeduct := order.FinalOrderValue

	var walletTransaction models.WalletTransaction

	walletTransaction.Amount = totalAmountToDeduct
	walletTransaction.TransactionType = "debit"
	walletTransaction.UserId = order.UserId
	walletTransaction.TransactionStatus = "success"

	if err := tx.Create(&walletTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	paymentTransaction.PaymentStatus = "success"
	order.OrderStatus = "placed"

	if err := tx.Save(&paymentTransaction).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	order.PaymentTransactionId = paymentTransaction.PaymentTransactionId

	if err := tx.Preload("User").Where("order_id=?", request.OrderID).Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	if err := tx.Where("user_id = ?", order.UserId).Delete(&cartItem).Error; err != nil {
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
		Where("user_id = ?", order.UserId).
		Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete cart items"})
	}

	var user models.User
	if err := tx.First(&user, "user_id = ?", order.UserId).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	NotifyOrderPlaced(order.OrderId)
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Payment Completed", "payment_transaction": paymentTransaction, "order": order, "wallet_transaction": walletTransaction, "response": fiber.Map{"order_status": "PAID"}})

}

func VerifyRazorpaySignature(orderID, paymentID, razorpaySignature, secret string) bool {
	// Concatenate order_id and payment_id as per Razorpay docs
	data := orderID + "|" + paymentID

	// Generate HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	generatedSignature := hex.EncodeToString(h.Sum(nil))

	// Compare signatures
	return generatedSignature == razorpaySignature
}

func RazorpayWebhookHandler(c *fiber.Ctx) error {

	body := c.Body()

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return c.Status(500).SendString("Failed to parse webhook body")
	}

	event := payload["event"].(string)

	if event == "payment.captured" {
		payment := payload["payload"].(map[string]interface{})["payment"].(map[string]interface{})["entity"].(map[string]interface{})

		orderID := payment["order_id"].(string)

		body, err := rz_client.Order.Fetch(orderID, nil, nil)
		tx := database.DB.Begin()

		var paymentTransaction models.PaymentTransaction

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if body["status"] != "paid" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Order is not paid  yet "})
		}

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		if err := tx.First(&paymentTransaction, "payment_order_id = ?", orderID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		// ✅ Update DB - mark order as paid using orderID/paymentID

		var order models.Order

		var cartItem models.CartItem

		// cf_order_id := c.Params("cf_order_id"

		if err := tx.Where("payment_order_id=?", orderID).First(&paymentTransaction).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}

		if paymentTransaction.PaymentStatus == "success" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Payment Already Completed"})
		}

		if err := tx.First(&order, "order_id=?", paymentTransaction.OrderId).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}

		if order.OrderStatus != "pending" && order.OrderStatus == "placed" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Order is already Verifed and Placed !", "response": fiber.Map{"order_status": "PAID"}})
		}

		totalAmountToDeduct := order.FinalOrderValue

		var walletTransaction models.WalletTransaction

		walletTransaction.Amount = totalAmountToDeduct
		walletTransaction.TransactionType = "debit"
		walletTransaction.UserId = order.UserId
		walletTransaction.TransactionStatus = "success"

		if err := tx.Create(&walletTransaction).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}

		paymentTransaction.PaymentStatus = "success"
		order.OrderStatus = "placed"

		if err := tx.Save(&paymentTransaction).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}

		order.PaymentTransactionId = paymentTransaction.PaymentTransactionId

		if err := tx.Preload("User").Where("order_id=?", paymentTransaction.OrderId).Save(&order).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
		}
		if err := tx.Where("user_id = ?", order.UserId).Delete(&cartItem).Error; err != nil {
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
			Where("user_id = ?", order.UserId).
			Delete(&models.CartItem{}).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete cart items"})
		}
		NotifyOrderPlaced(order.OrderId)

		if err := tx.Commit().Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
		}

	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Payment  Completed succesfully",
	})
}
