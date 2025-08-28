package controllers

import (
	"backend-hostego/config"
	"backend-hostego/database"
	"backend-hostego/logs"
	"backend-hostego/models"
	natsclient "backend-hostego/nats"
	"crypto/hmac"
	"crypto/sha256"
	"runtime/debug"
	"strconv"
	"time"

	"backend-hostego/middlewares"

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
	userId, err := middlewares.SafeUserIDExtractor(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user authentication: " + err.Error()})
	}
	if userId == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	type OrderRequest struct {
		OrderID int `json:"order_id"` //Only accept Order ID, not amount
	}
	var request OrderRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	var result struct {
		Order              models.Order
		Wallet             models.Wallet
		WalletTransaction  models.WalletTransaction
		PaymentTransaction models.PaymentTransaction
		OrderItems         []models.CartItem
	}

	err = database.SafeTransactionWithCleanup(func(tx *gorm.DB) error {
		// Fetch order
		if err := tx.First(&result.Order, "order_id=?", request.OrderID).Error; err != nil {
			return err
		}

		// Fetch wallet
		if err := tx.First(&result.Wallet, "user_id=?", userId).Error; err != nil {
			return err
		}

		// Check wallet balance
		if result.Wallet.Balance < result.Order.FinalOrderValue {
			return fmt.Errorf("wallet balance insufficient to complete payment")
		}

		// Deduct amount from wallet
		totalAmountToDeduct := result.Order.FinalOrderValue
		result.Wallet.Balance -= totalAmountToDeduct

		// Create wallet transaction
		result.WalletTransaction.Amount = totalAmountToDeduct
		result.WalletTransaction.TransactionType = "debit"
		result.WalletTransaction.UserId = userId
		result.WalletTransaction.TransactionStatus = "success"

		if err := tx.Create(&result.WalletTransaction).Error; err != nil {
			return err
		}

		// Create payment transaction
		result.PaymentTransaction.OrderId = request.OrderID
		result.PaymentTransaction.UserId = userId
		result.PaymentTransaction.Amount = totalAmountToDeduct
		result.PaymentTransaction.PaymentStatus = "success"
		result.PaymentTransaction.PaymentMethod = "wallet"

		if err := tx.Create(&result.PaymentTransaction).Error; err != nil {
			return err
		}

		// Update order
		result.Order.PaymentTransactionId = result.PaymentTransaction.PaymentTransactionId
		result.Order.OrderStatus = "placed"

		// Update wallet
		if err := tx.Where("user_id=?", userId).Save(&result.Wallet).Error; err != nil {
			return err
		}

		// Parse order items
		if err := json.Unmarshal(result.Order.OrderItems, &result.OrderItems); err != nil {
			return fmt.Errorf("failed to parse order items: %v", err)
		}

		// Process each order item
		for _, item := range result.OrderItems {
			result.Order.RestaurantPayableAmount += item.ActualSubTotal

			orderItem := models.OrderItem{
				OrderId:        result.Order.OrderId,
				ProductId:      item.ProductId,
				Quantity:       item.Quantity,
				SubTotal:       item.SubTotal,
				UserId:         result.Order.UserId,
				ActualSubTotal: item.ActualSubTotal,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				return err
			}

			// Update product stock
			if err := tx.Model(&models.Product{}).
				Where("product_id = ?", item.ProductId).
				Update("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity)).Error; err != nil {
				return err
			}
		}

		// Save updated order
		if err := tx.Preload("User").Where("order_id=?", request.OrderID).Save(&result.Order).Error; err != nil {
			return err
		}

		// Delete cart items
		if err := tx.Where("user_id = ?", userId).Delete(&models.CartItem{}).Error; err != nil {
			return fmt.Errorf("failed to delete cart items: %v", err)
		}

		return nil
	}, "InitiatePayment")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Send notification outside transaction
	NotifyOrderPlaced(result.Order.OrderId)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":             "Payment Completed",
		"payment_transaction": result.PaymentTransaction,
		"order":               result.Order,
		"wallet_transaction":  result.WalletTransaction,
	})
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
	current_user_id, err := middlewares.SafeUserIDExtractor(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user authentication: " + err.Error()})
	}
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

	var result struct {
		Order             models.Order
		Wallet            models.Wallet
		WalletTransaction models.WalletTransaction
		OrderItems        []models.CartItem
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

		// Update order status
		result.Order.OrderStatus = models.OrderStatusType(models.CanceledOrderStatus)
		result.Order.Refunded = true
		result.Order.RefundedAt = time.Now()
		result.Order.RefundInitiator = current_user_id
		result.Order.DeliveryPartnerId = 0
		result.Order.DeliveryPartner = nil

		if err := tx.Save(&result.Order).Error; err != nil {
			return err
		}

		// Create wallet transaction for refund
		result.WalletTransaction.Amount = result.Order.FinalOrderValue
		result.WalletTransaction.TransactionType = models.TransactionCustomType(models.TransactionRefund)
		result.WalletTransaction.TransactionStatus = models.TransactionStatusType(models.TransactionSuccess)
		result.WalletTransaction.UserId = result.Order.UserId

		if err := tx.Create(&result.WalletTransaction).Error; err != nil {
			return err
		}

		// Update wallet balance
		result.Wallet.Balance += result.Order.FinalOrderValue
		if err := tx.Where("user_id=?", result.Order.UserId).Save(&result.Wallet).Error; err != nil {
			return err
		}

		// Parse order items
		if err := json.Unmarshal(result.Order.OrderItems, &result.OrderItems); err != nil {
			return fmt.Errorf("failed to parse order items: %v", err)
		}

		// Process refund for each item
		for _, item := range result.OrderItems {
			// Restore product stock
			if err := tx.Model(&models.Product{}).Where("product_id = ?", item.ProductId).
				Update("stock_quantity", gorm.Expr("stock_quantity + ?", item.Quantity)).Error; err != nil {
				return err
			}

			// Delete order item
			if err := tx.Where("order_id = ? AND product_id = ?", result.Order.OrderId, item.ProductId).
				Delete(&models.OrderItem{}).Error; err != nil {
				return err
			}
		}

		return nil
	}, "InitiateRefundPayment")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":            "Refund Completed",
		"wallet_transaction": result.WalletTransaction,
		"wallet":             result.Wallet,
	})
}

func InitateCashfreePaymentOrder(c *fiber.Ctx) error {
	user_id := c.Locals("user_id").(int)
	if user_id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Unauthorized"})
	}

	var orderRequest struct {
		OrderId int `json:"order_id"`
	}
	if err := c.BodyParser(&orderRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var result struct {
		User               models.User
		Order              models.Order
		PaymentTransaction models.PaymentTransaction
		CashfreeResponse   map[string]interface{}
	}

	err := database.SafeTransactionWithCleanup(func(tx *gorm.DB) error {
		// Fetch user
		if err := tx.First(&result.User, "user_id = ?", user_id).Error; err != nil {
			return err
		}

		// Fetch order
		if err := tx.First(&result.Order, "order_id = ?", orderRequest.OrderId).Error; err != nil {
			return err
		}

		// Prepare Cashfree request body
		body := map[string]interface{}{
			"order_amount":   result.Order.FinalOrderValue,
			"order_currency": "INR",
			"customer_details": map[string]interface{}{
				"customer_id":    strconv.Itoa(result.User.UserId),
				"customer_phone": result.User.MobileNumber,
				"customer_email": result.User.Email,
				"customer_name":  result.User.FirstName,
			},
		}

		// Make Cashfree API call
		restyClient := resty.New()
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
			return err
		}

		// Parse Cashfree response
		if err := json.Unmarshal(resp.Body(), &result.CashfreeResponse); err != nil {
			return fmt.Errorf("invalid response from Cashfree: %v", err)
		}

		// Create payment transaction
		result.PaymentTransaction.OrderId = orderRequest.OrderId
		result.PaymentTransaction.UserId = user_id
		result.PaymentTransaction.Amount = result.Order.FinalOrderValue
		result.PaymentTransaction.PaymentStatus = "pending"
		result.PaymentTransaction.PaymentMethod = "UPI"
		result.PaymentTransaction.PaymentOrderId = result.CashfreeResponse["order_id"].(string)

		if err := tx.Create(&result.PaymentTransaction).Error; err != nil {
			return err
		}

		// Update order with payment transaction ID
		result.Order.PaymentTransactionId = result.PaymentTransaction.PaymentTransactionId
		if err := tx.Save(&result.Order).Error; err != nil {
			return err
		}

		return nil
	}, "InitateCashfreePaymentOrder")

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(result.CashfreeResponse)
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
			OrderId:        order.OrderId,
			ProductId:      item.ProductId,
			Quantity:       item.Quantity,
			SubTotal:       item.SubTotal,
			UserId:         order.UserId,
			ActualSubTotal: item.ActualSubTotal,
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
		order.RestaurantPayableAmount += item.ActualSubTotal
		orderItem := models.OrderItem{
			OrderId:        order.OrderId,
			ProductId:      item.ProductId,
			Quantity:       item.Quantity,
			SubTotal:       item.SubTotal,
			ActualSubTotal: item.ActualSubTotal,
			UserId:         order.UserId,
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
	if err := tx.Preload("User").Where("order_id=?", request.OrderID).Save(&order).Error; err != nil {
		tx.Rollback()
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

// RazorpayWebhookHandler receives Razorpay webhook and quickly responds
// RazorpayWebhookHandler receives Razorpay webhook and quickly responds
func RazorpayWebhookHandler(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in RazorpayWebhookHandler:", r)
		}
	}()

	body := c.Body()

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Println("Failed to parse webhook:", err)
		return c.Status(500).SendString("Invalid payload")
	}

	event, ok := payload["event"].(string)
	if !ok || event != "payment.captured" {
		log.Println("Unsupported or missing event")
		return c.SendStatus(fiber.StatusOK)
	}

	// ✅ Launch goroutine safely with enhanced recovery
	go func() {
		defer func() {
			if r := recover(); r != nil {
				additionalData := map[string]interface{}{
					"payload":  payload,
					"function": "ProcessPaymentCaptured",
				}
				logs.LogCrash(r, "Payment Processing Goroutine", additionalData)
			}
		}()
		ProcessPaymentCaptured(payload)
	}()

	return c.Status(200).JSON(fiber.Map{
		"message": "Payment captured event received",
	})
}
func ProcessPaymentCaptured(payload map[string]interface{}) {
	defer func() {
		if r := recover(); r != nil {
			additionalData := map[string]interface{}{
				"payload": payload,
				"stack":   string(debug.Stack()),
			}
			logs.LogCrash(r, "ProcessPaymentCaptured", additionalData)
		}
	}()

	// Safely extract nested fields
	payloadData, ok := payload["payload"].(map[string]interface{})
	if !ok {
		log.Println("Invalid payload: missing 'payload'")
		return
	}

	paymentData, ok := payloadData["payment"].(map[string]interface{})
	if !ok {
		log.Println("Invalid payload: missing 'payment'")
		return
	}

	entity, ok := paymentData["entity"].(map[string]interface{})
	if !ok {
		log.Println("Invalid payload: missing 'entity'")
		return
	}

	orderID, ok := entity["order_id"].(string)
	if !ok {
		log.Println("Missing order_id in entity")
		return
	}

	log.Println("Processing payment for order ID:", orderID)

	// Fetch order from Razorpay
	body, err := rz_client.Order.Fetch(orderID, nil, nil)
	if err != nil {
		log.Println("Error fetching Razorpay order:", err)
		return
	}

	status, ok := body["status"].(string)
	if !ok || status != "paid" {
		log.Println("Order not marked as paid by Razorpay, status:", status)
		return
	}

	var result struct {
		PaymentTransaction models.PaymentTransaction
		Order              models.Order
		WalletTransaction  models.WalletTransaction
		OrderItems         []models.CartItem
	}

	err = database.SafeTransactionWithCleanup(func(tx *gorm.DB) error {
		// Fetch payment transaction
		if err := tx.First(&result.PaymentTransaction, "payment_order_id = ?", orderID).Error; err != nil {
			return fmt.Errorf("payment transaction fetch error: %v", err)
		}

		if result.PaymentTransaction.PaymentStatus == "success" {
			return fmt.Errorf("payment already completed for orderID: %s", orderID)
		}

		// Fetch order
		if err := tx.First(&result.Order, "order_id = ?", result.PaymentTransaction.OrderId).Error; err != nil {
			return fmt.Errorf("order fetch error: %v", err)
		}

		if result.Order.OrderStatus != "pending" {
			return fmt.Errorf("order already placed, skipping update. OrderID: %d", result.Order.OrderId)
		}

		// Create wallet transaction
		totalAmountToDeduct := result.Order.FinalOrderValue
		result.WalletTransaction = models.WalletTransaction{
			Amount:            totalAmountToDeduct,
			TransactionType:   "debit",
			UserId:            result.Order.UserId,
			TransactionStatus: "success",
		}

		if err := tx.Create(&result.WalletTransaction).Error; err != nil {
			return fmt.Errorf("wallet transaction creation failed: %v", err)
		}

		// Update payment & order status
		result.PaymentTransaction.PaymentStatus = "success"
		result.Order.OrderStatus = "placed"
		result.Order.PaymentTransactionId = result.PaymentTransaction.PaymentTransactionId

		if err := tx.Save(&result.PaymentTransaction).Error; err != nil {
			return fmt.Errorf("saving payment transaction failed: %v", err)
		}

		// Delete cart items
		if err := tx.Where("user_id = ?", result.Order.UserId).Delete(&models.CartItem{}).Error; err != nil {
			return fmt.Errorf("cart item deletion failed: %v", err)
		}

		// Parse order items
		if err := json.Unmarshal(result.Order.OrderItems, &result.OrderItems); err != nil {
			return fmt.Errorf("failed to unmarshal OrderItems JSON: %v", err)
		}

		// Process each order item
		for _, item := range result.OrderItems {
			result.Order.RestaurantPayableAmount += item.ActualSubTotal

			orderItem := models.OrderItem{
				OrderId:        result.Order.OrderId,
				ProductId:      item.ProductId,
				Quantity:       item.Quantity,
				SubTotal:       item.SubTotal,
				UserId:         result.Order.UserId,
				ActualSubTotal: item.ActualSubTotal,
			}

			if err := tx.Create(&orderItem).Error; err != nil {
				return fmt.Errorf("failed to create OrderItem: %v", err)
			}

			if err := tx.Model(&models.Product{}).
				Where("product_id = ?", item.ProductId).
				Update("stock_quantity", gorm.Expr("stock_quantity - ?", item.Quantity)).Error; err != nil {
				return fmt.Errorf("failed to update product stock: %v", err)
			}
		}

		// Save updated order
		if err := tx.Preload("User").Save(&result.Order).Error; err != nil {
			return fmt.Errorf("failed to save updated order: %v", err)
		}

		return nil
	}, "ProcessPaymentCaptured")

	if err != nil {
		log.Printf("❌ ProcessPaymentCaptured failed for orderID %s: %v", orderID, err)
		return
	}

	// Send notification outside transaction
	if err := NotifyOrderPlaced(result.Order.OrderId); err != nil {
		log.Printf("⚠️ Notification sending failed for order %d: %v", result.Order.OrderId, err)
	}

	log.Printf("✅ Payment processing completed successfully for orderID: %s", orderID)
}

func InitiateWalletTopup(c *fiber.Ctx) error {

	type WalletTopupRequest struct {
		Amount int `json:"amount"`
	}

	request := new(WalletTopupRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	var user models.User
	if err := database.DB.First(&user, "user_id = ?", c.Locals("user_id")).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	walletTransaction := models.WalletTransaction{
		Amount:            float64(request.Amount),
		TransactionType:   "credit",
		UserId:            user.UserId,
		TransactionStatus: "success",
		PaymentOrderId:    "wallet_topup",
	}

	if err := database.DB.Create(&walletTransaction).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Wallet topup successful", "wallet_transaction": walletTransaction})
}
