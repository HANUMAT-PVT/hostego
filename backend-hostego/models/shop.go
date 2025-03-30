package models

import "time"

type Shop struct {
	ShopId          int          `gorm:"type:int;primaryKey;autoIncrement:true" json:"shop_id"`
	ShopName        string       `gorm:"type:varchar(255);" json:"shop_name"`
	ShopImg         string       `gorm:"type:varchar(255);" json:"shop_img"`
	Address         string       `gorm:"type:varchar(255);" json:"address"`
	PreparationTime string       `gorm:"type:varchar(20);"  json:"preparation_time"`
	FoodCategory    FoodCategory `gorm:"embedded" json:"food_category"`
	ShopStatus      int          `gorm:"type:int;default:1;index;" json:"shop_status"`
	CreatedAt       time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
}
