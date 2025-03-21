package models

import "time"

type DeliveryPartnerWallet struct {
	DeliveryPartnerWalletId int             `gorm:"type:int;primaryKey;unique;not null;autoIncrement:true" json:"delivery_partner_wallet_id"`
	DeliveryPartnerId       string          `gorm:"type:uuid;not null;index" json:"delivery_partner_id"` // Foreign Key

	Balance                 float64         `gorm:"type:double precision;default:0" json:"balance"`
	CreatedAt               time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	

}
