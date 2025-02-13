package models

import "github.com/google/uuid"

type Shop struct {
	ShopId          uuid.UUID    `gorm:"type:uuid;primaryKey;not null;unique;default:gen_random_uuid();" json:"shop_id"`
	ShopName        string       `gorm:"type:varchar(255);" json:"shop_name"`
	ShopImg         string       `gorm:"type:varchar(255);" json:"shop_img"`
	Address         string       `gorm:"type:varchar(255);" json:"address"`
	PreparationTime string       `gorm:"type:varchar(20);"  json:"preparation_time"`
	FoodCategory    FoodCategory `gorm:"embedded" json:"food_category"`
}
