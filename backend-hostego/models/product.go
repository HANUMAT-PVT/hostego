package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
	ProductId       string       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"product_id"`
	ProductName     string       `gorm:"type:varchar(255);not null" json:"product_name"`
	FoodCategory    FoodCategory `gorm:"embedded" json:"food_category"`
	FoodPrice       float64      `gorm:"type:double precision;not null" json:"food_price"`
	Availability    int          `gorm:"type:int;not null;default:1" json:"availability"`
	ProductImg      string       `gorm:"type:varchar(255);not null" json:"product_img_url"`
	Description     string       `gorm:"type:varchar(255);not null" json:"description"`
	Discount        Discount     `gorm:"embedded" json:"discount"`
	CreatedAt       time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	PreparationTime string       `gorm:"type:varchar(255);not null" json:"preparation_time"`
	ShopId          string       `gorm:"type:uuid;index;not null" json:"shop_id"`
	Shop            Shop         `gorm:"foreignKey:shop_id;references:Shop" json:"shop"`
}

// BeforeCreate hook to generate UUID before inserting a new product
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	p.ProductId = uuid.New().String()
	return nil
}
