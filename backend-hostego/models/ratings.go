package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Rating struct {
	RatingId   string    `gorm:"type:uuid;primaryKey;not null;unique;default:gen_random_uuid();" json:"rating_id"`
	ProductID  string    `gorm:"type:uuid;index;not null;" json:"product_id"`
	UserID     string    `gorm:"type:uuid;not null;" json:"user_id"`
	Rating     int       `gorm:"type:int;not null;check:rating BETWEEN 1 AND 5" json:"rating"`
	ReviewText string    `gorm:"type:text" json:"review_text"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	Product    Product   `gorm:"type:uuid;foreignKey:product_id;refrences:Product" json:"product"`
	User       User      `gorm:"type:uuid;foreignKey:user_id;refrences:User" json:"user"`
}

func (r *Rating) BeforeCreate(trx *gorm.DB) (err error){
	r.RatingId=uuid.New().String();
	return nil;
}


