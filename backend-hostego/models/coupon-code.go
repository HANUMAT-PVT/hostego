package models

import "time"

type CouponCode struct {
	Id               int       `gorm:"type:int;primaryKey;autoIncrement:true" json:"id"`
	Code             string    `gorm:"type:string;not null" json:"code"`
	Discount         float64   `gorm:"type:double precision;not null" json:"discount"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ShopId           int       `gorm:"type:int" json:"shop_id"`
	ProductId        int       `gorm:"type:int" json:"product_id"`
	IsActive         bool      `gorm:"type:bool;default:true" json:"is_active"`
	InfluencerUserId int       `gorm:"type:int;" json:"influencer_user_id"`
	MinOrderAmount   float64   `gorm:"type:double precision;default:0" json:"min_order_amount"`
}
