package models

import (
	"time"
)

type Rating struct {
	RatingId   int       `gorm:"type:int;primaryKey;autoIncrement:true" json:"rating_id"`
	OrderID    int       `gorm:"type:int;index;not null;" json:"order_id"`
	ProductID  int       `gorm:"type:int;index;not null;" json:"product_id"`
	UserID     int       `gorm:"type:int;not null;" json:"user_id"`
	Rating     float64   `gorm:"type:double precision;not null;" json:"rating"`
	ReviewText string    `gorm:"type:text" json:"review_text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
