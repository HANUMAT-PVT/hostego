package models

import (
	"time"

	"gorm.io/gorm"
)

type Rating struct {
	RatingId   int    `gorm:"type:int;primaryKey;autoIncrement:true" json:"rating_id"`
	ProductID  int    `gorm:"type:int;index;not null;" json:"product_id"`
	UserID     int    `gorm:"type:int;not null;" json:"user_id"`
	Rating     int       `gorm:"type:int;not null;check:rating BETWEEN 1 AND 5" json:"rating"`
	ReviewText string    `gorm:"type:text" json:"review_text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	Product    Product   `gorm:"type:uuid;foreignKey:product_id;references:Product" json:"product"`
	User       User      `gorm:"type:uuid;foreignKey:user_id;references:User" json:"user"`
}

func (r *Rating) BeforeCreate(trx *gorm.DB) (err error) {
	r.RatingId = r.ProductID
	return nil
}
