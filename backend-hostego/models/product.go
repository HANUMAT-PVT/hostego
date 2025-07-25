package models

import (
	"time"

	"gorm.io/datatypes"
)

// Food Category Structure
type FoodCategory struct {
	IsVeg    int `gorm:"type:int;not null;default:0" json:"is_veg"`
	IsCooked int `gorm:"type:int;not null;default:0" json:"is_cooked"`
}

// Discount Structure
type Discount struct {
	IsAvailable int `gorm:"type:int;not null;" json:"is_available"`
	Percentage  int `gorm:"type:int;not null;" json:"percentage"`
}

// Product Model
type Product struct {
	ProductId       int            `gorm:"type:int;primaryKey;unique;autoIncrement:true" json:"product_id"`
	ProductName     string         `gorm:"type:varchar(255);index;not null" json:"product_name"`
	FoodCategory    FoodCategory   `gorm:"embedded" json:"food_category"`
	FoodPrice       float64        `gorm:"type:double precision;not null" json:"food_price"`
	Availability    int            `gorm:"type:int;not null;default:1" json:"availability"`
	ProductImgUrl   string         `gorm:"type:varchar(255);not null" json:"product_img_url"`
	Description     string         `gorm:"type:varchar(255);index;not null" json:"description"`
	Discount        Discount       `gorm:"embedded" json:"discount"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	PreparationTime string         `gorm:"type:varchar(255);" json:"preparation_time"`
	ShopId          int            `gorm:"type:int;not null;index" json:"shop_id"`
	Shop            Shop           `gorm:"foreignKey:ShopId;references:ShopId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"shop"`
	Tags            datatypes.JSON `gorm:"type:jsonB;" json:"tags"`
	StockQuantity   int            `gorm:"type:int;not null;default:0" json:"stock_quantity"`
	Weight          string         `gorm:"type:varchar(255);" json:"weight"`
	SellingPrice    float64        `gorm:"type:double precision" json:"selling_price"`
	Category        string         `gorm:"type:varchar(255);" json:"category"`
	AverageRating   float64        `gorm:"type:double precision;default:0" json:"average_rating"`
	TotalRatings    int            `gorm:"type:int;default:0" json:"total_ratings"`
}
