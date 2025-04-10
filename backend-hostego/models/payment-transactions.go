package models

import (
	"time"
)

//	Payment  PaymentTransactions Schema{
//	  "payment_PaymentTransaction_id": "txn_123456",
//	  "order_id": "ORD_123456",
//	  "user_id": "user_987654321",
//	  "amount": 250.00,
//	  "payment_method": "wallet",
//	  "payment_status": "success",
//	  "payment_timestamp": "2025-02-09T12:15:00Z"
//	}
type PaymentTransactionStatusType string

const (
	PaymentTransactionPending  PaymentTransactionStatusType = "pending"
	PaymentTransactionSuccess  PaymentTransactionStatusType = "success"
	PaymentTransactionFailed   PaymentTransactionStatusType = "failed"
	PaymentTransactionRefunded PaymentTransactionStatusType = "refunded"
	PaymentTransactionCanceled PaymentTransactionStatusType = "canceled"
)

type PaymentTransaction struct {
	PaymentTransactionId int                          `gorm:"type:int;primaryKey;not null;unique;autoIncrement:true" json:"payment_PaymentTransaction_id"`
	OrderId              int                          `gorm:"type:int;not null;" json:"order_id"`
	UserId               int                          `gorm:"type:int;not null;" json:"user_id"`
	Amount               float64                      `gorm:"type:double precision;not null;" json:"amount"`
	PaymentMethod        string                       `gorm:"type:varchar(20);not null;" json:"payment_method"`
	PaymentStatus        PaymentTransactionStatusType `gorm:"type:varchar(20);not null" json:"payment_transaction_status"`
	CreatedAt            time.Time                    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time                    `gorm:"autoUpdateTime" json:"updated_at"`
	User                 User                         `gorm:"foreignKey:UserId;references:UserId;" json:"user"`
	PaymentOrderId       string                       `gorm:"type:string;" json:"payment_order_id"`
}
