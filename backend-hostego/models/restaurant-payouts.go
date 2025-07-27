package models

import "time"

type RestaurantPayout struct {
	PayoutID      uint       `gorm:"primaryKey;autoIncrement" json:"payout_id"`
	OrderID       uint       `gorm:"not null;index" json:"order_id"`
	ShopID        uint       `gorm:"not null;index" json:"shop_id"`
	Amount        float64    `gorm:"not null" json:"amount"`
	Status        string     `gorm:"type:varchar(20);default:'pending'" json:"status"` // pending, paid, failed
	PaidAt        *time.Time `json:"paid_at,omitempty"`
	PaymentRef    string     `gorm:"type:varchar(100)" json:"payment_ref"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	Note          string     `gorm:"type:varchar(255)" json:"note"`
	PaymentMethod string     `gorm:"type:varchar(255)" json:"payment_method"`
}
