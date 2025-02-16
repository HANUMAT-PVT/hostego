package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// Food Category Structure
type FoodCategory struct {
	IsVeg    int `gorm:"type:int;not null;" json:"is_veg"`
	IsCooked int `gorm:"type:int;not null;" json:"is_cooked"`
}

// Discount Structure
type Discount struct {
	IsAvailable int `gorm:"type:int;not null;" json:"is_available"`
	Percentage  int `gorm:"type:int;not null;" json:"percentage"`
}

// Product Model
type Product struct {
	ProductId       string         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"product_id"`
	ProductName     string         `gorm:"type:varchar(255);index;not null" json:"product_name"`
	FoodCategory    FoodCategory   `gorm:"embedded" json:"food_category"`
	FoodPrice       float64        `gorm:"type:double precision;not null" json:"food_price"`
	Availability    int            `gorm:"type:int;not null;default:1" json:"availability"`
	ProductImg      string         `gorm:"type:varchar(255);not null" json:"product_img_url"`
	Description     string         `gorm:"type:varchar(255);index;not null" json:"description"`
	Discount        Discount       `gorm:"embedded" json:"discount"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	PreparationTime string         `gorm:"type:varchar(255);not null" json:"preparation_time"`
	ShopId          uuid.UUID      `gorm:"type:uuid;not null;index" json:"shop_id"`
	Shop            Shop           `gorm:"foreignKey:ShopId;references:ShopId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"shop"`
	Tags            datatypes.JSON `gorm:"type:jsonB;" json:"tags"`
}
