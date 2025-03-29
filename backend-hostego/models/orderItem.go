package models

import (
	"time"
)

type OrderItem struct {
	OrderItemId int       `gorm:"type:int;primaryKey;autoIncrement:true" json:"order_item_id"`
	OrderId     int       `gorm:"type:int;not null;index" json:"order_id"`
	ProductId   int       `gorm:"type:int;not null" json:"product_id"`
	UserId      int       `gorm:"type:int;not null;index" json:"user_id"`
	Quantity    int       `gorm:"type:int;not null" json:"quantity"`
	SubTotal    float64   `gorm:"type:double precision;not null" json:"sub_total"`
	Product     Product   `gorm:"foreignKey:ProductId;references:ProductId" json:"product"`
	Order       Order     `gorm:"foreignKey:OrderId;references:OrderId" json:"order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
