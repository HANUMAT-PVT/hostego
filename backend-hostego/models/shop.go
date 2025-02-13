package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Shop struct {
	ShopId          string       `gorm:"type:uuid;primaryKey;not null;unique;default:gen_random_uuid();" json:"shop_id"`
	ShopName        string       `gorm:"type:varchar(255);" json:"shop_name"`
	ShopImg         string       `gorm:"type:varchar(255);" json:"shop_img"`
	Address         string       `gorm:"type:varchar(255);" json:"address"`
	PreparationTime string       `gorm:"type:varchar(20);"  json:"preparation_time"`
	FoodCategory    FoodCategory `gorm:"embedded" json:"food_category"`
}

// BeforeCreate hook to generate UUID before inserting a new transaction
func (s *Shop) BeforeCreate(tx *gorm.DB) (err error) {
	s.ShopId = uuid.New().String()
	return nil
}
