package models

import "time"

type Wallet struct {
	WalletId  int       `gorm:"type:int;primaryKey;unique;not null;autoIncrement:true" json:"wallet_id"`
	UserId    int       `gorm:"type:int;not null;index" json:"user_id"` // Foreign Key
	User      User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;references:UserId" json:"user"`
	Balance   float64   `gorm:"type:double precision;default:0" json:"balance"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
