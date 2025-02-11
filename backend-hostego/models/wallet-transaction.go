package models

import (
	"time"
)

type PaymentMethod struct {
	Type                string `gorm:"type:varchar(50);not null" json:"type"` // e.g., "UPI", "Card", "NetBanking"
	UniqueTransactionID string `gorm:"type:varchar(100);not null;unique" json:"unique_transaction_id"`
}

type WalletTransaction struct {
	TransactionId     string        `gorm:"type:primaryKey;uuid;not null;unique;default:gen_random_uuid()" json:"transaction_id"`
	UserId            string        `gorm:"type:uuid;index;not null;" json:"user_id"`
	User              User          `gorm:"type:foreignKey:User;" json:"user"`
	Amount            float64       `gorm:"type:float64;not null" json:"amount"`
	TransactionType   string        `gorm:"type:varchaar(20);not null;" json:"transaction_type"`
	CreatedAt         time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	TransactionStatus string        `gorm:"type:varchaar(20);not null;" json:"transaction_status"`
	PaymentMethod     PaymentMethod `gorm:"embedded" json:"payment_method"`
}
