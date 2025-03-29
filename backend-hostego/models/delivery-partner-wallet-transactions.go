package models

import (
	"time"
)

// DeliveryWalletTransaction represents a user's wallet transaction
type DeliveryPartnerWalletTransaction struct {
	TransactionID     int                   `gorm:"type:int;primaryKey;unique;not null;autoIncrement:true" json:"transaction_id"`
	DeliveryPartnerId int                   `gorm:"type:int;not null;index" json:"delivery_partner_id"`
	
	Amount            float64               `gorm:"type:double precision;not null" json:"amount"`
	TransactionType   TransactionCustomType `gorm:"type:varchar(20);not null" json:"transaction_type"`
	CreatedAt         time.Time             `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time             `gorm:"autoUpdateTime" json:"updated_at"`
	TransactionStatus TransactionStatusType `gorm:"type:varchar(20);not null" json:"transaction_status"`
	PaymentMethod     PaymentMethod         `gorm:"embedded" json:"payment_method"`

}
