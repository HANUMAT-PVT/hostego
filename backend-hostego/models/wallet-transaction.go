package models

import (
	"time"
)

// PaymentMethod represents the method used for the transaction
type PaymentMethod struct {
	Type                    string `gorm:"type:varchar(50);not null" json:"type"` // e.g., "UPI", "Card", "NetBanking"
	UniqueTransactionID     string `gorm:"type:varchar(100);not null;unique" json:"unique_transaction_id"`
	PaymentScreenShotImgUrl string `gorm:"type:varchar(100);not null;unique" json:"payment_screenshot_img_url"`
	PaymentVerifiedByAdmin  string `gorm:"type:uuid" json:"payment_verified_by_admin"`
}
type TransactionStatusType string

const (
	TransactionPending  TransactionStatusType = "pending"
	TransactionSuccess  TransactionStatusType = "success"
	TransactionFailed   TransactionStatusType = "failed"
	TransactionRefunded TransactionStatusType = "refunded"
	TransactionCanceled TransactionStatusType = "canceled"
)

type TransactionCustomType string

const (
	TransactionDebit  string = "debit"
	TransactionCredit string = "credit"
	TransactionRefund string = "refund"
)

// WalletTransaction represents a user's wallet transaction
type WalletTransaction struct {
	TransactionID     string                `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"transaction_id"`
	UserID            string                `gorm:"type:uuid;not null;index" json:"user_id"`
	User              User                  `gorm:"foreignKey:user_id;references:User" json:"user"`
	Amount            float64               `gorm:"type:double precision;not null" json:"amount"`
	TransactionType   TransactionCustomType `gorm:"type:varchar(20);not null" json:"transaction_type"`
	CreatedAt         time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	TransactionStatus TransactionStatusType `gorm:"type:varchar(20);not null" json:"transaction_status"`
	PaymentMethod     PaymentMethod         `gorm:"embedded" json:"payment_method"`
}
