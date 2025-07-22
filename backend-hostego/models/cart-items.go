package models

import (
	"time"
)

type CartItem struct {
	CartItemId int       `gorm:"type:int;primaryKey;unique;not null;autoIncrement:true" json:"cart_item_id"`
	ProductId  int       `gorm:"type:int;not null;index" json:"product_id"`
	Quantity   int       `gorm:"type:int;default:1;not null;" json:"quantity"`
	SubTotal   float64   `gorm:"type:double precision;not null;" json:"sub_total"`
	UserId     int       `gorm:"type:int;not null;index;" json:"user_id"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	// Foreign Key Relation
	ProductItem Product `gorm:"foreignKey:ProductId;references:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product_item"`
}
