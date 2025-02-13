package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentMethod represents the method used for the transaction
type PaymentMethod struct {
	Type                string `gorm:"type:varchar(50);not null" json:"type"` // e.g., "UPI", "Card", "NetBanking"
	UniqueTransactionID string `gorm:"type:varchar(100);not null;unique" json:"unique_transaction_id"`
}

// WalletTransaction represents a user's wallet transaction
type WalletTransaction struct {
	TransactionID     string        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"transaction_id"`
	UserID            string        `gorm:"type:uuid;not null;index" json:"user_id"`
	User              User          `gorm:"foreignKey:user_id;references:User" json:"user"`
	Amount            float64       `gorm:"type:double precision;not null" json:"amount"`
	TransactionType   string        `gorm:"type:varchar(20);not null" json:"transaction_type"`
	CreatedAt         time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	TransactionStatus string        `gorm:"type:varchar(20);not null" json:"transaction_status"`
	PaymentMethod     PaymentMethod `gorm:"embedded" json:"payment_method"`
}

// BeforeCreate hook to generate UUID before inserting a new transaction
func (w *WalletTransaction) BeforeCreate(tx *gorm.DB) (err error) {
	w.TransactionID = uuid.New().String()
	return nil
}
