package models

// Payment  Transactions Schema{
//   "payment_transaction_id": "txn_123456",
//   "order_id": "ORD_123456",
//   "user_id": "user_987654321",
//   "amount": 250.00,
//   "payment_method": "wallet",
//   "payment_status": "success",
//   "payment_timestamp": "2025-02-09T12:15:00Z"
// }

type PaymentTransaction struct {
	PaymentTransactionId string `gorm:"type:primaryKey;not null;unique;uuid;default:gen_random_uuid()" json:"payment_transaction_id"`
	OrderId              string `gorm:"type:foreignKey:Order;not null;" json:"order_id"`

}
