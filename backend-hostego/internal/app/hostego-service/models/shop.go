package models

import (
	"time"

	"github.com/google/uuid"
)

type Shop struct {
	ShopId             uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"shop_id"`
	UserId             int          `gorm:"not null;index" json:"user_id"`
	ShopName           string       `gorm:"type:varchar(255);not null" json:"shop_name"`
	ShopImg            []string     `gorm:"type:text[]" json:"shop_img"`
	Address            string       `gorm:"type:text;not null" json:"address"`
	AvgPreparationTime string       `gorm:"type:varchar(20);not null" json:"avg_preparation_time"`
	ShopEnabled        bool         `gorm:"type:boolean;default:true" json:"shop_enabled"`
	FoodCategory       FoodCategory `gorm:"embedded" json:"food_category"`
	CreatedAt          time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}
