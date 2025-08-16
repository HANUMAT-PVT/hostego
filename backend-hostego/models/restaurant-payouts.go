package models

import "time"

type RestaurantPayout struct {
	PayoutID      int        `gorm:"type:int;primaryKey;autoIncrement:true" json:"payout_id"`
	ShopID        int        `gorm:"index;not null" json:"shop_id"`
	Shop          Shop       `gorm:"foreignKey:ShopID;references:ShopId" json:"shop"`
	TotalAmount   float64    `gorm:"type:double precision;not null" json:"total_amount"`
	Status        string     `gorm:"type:varchar(50);default:'pending'" json:"status"` // pending, processing, paid, failed
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	PaidAt        *time.Time `gorm:"type:timestamp" json:"paid_at"`
	PaymentRefID  string     `gorm:"type:varchar(255)" json:"payment_ref_id"`
	PaymentMethod string     `gorm:"type:varchar(255)" json:"payment_method"` // reference from Razorpay/Stripe/Bank etc.
}
