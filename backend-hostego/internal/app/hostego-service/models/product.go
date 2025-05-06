package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type FoodCategory struct {
	IsVeg    bool `json:"is_veg"`
	IsNonVeg bool `json:"is_non_veg"`
}

type Discount struct {
	IsAvailable bool   `gorm:"not null;default:false" json:"is_available"`
	Percentage  string `gorm:"type:varchar(10);not null;default:'0'" json:"percentage"`
}

type Product struct {
	Id                 uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"product_id"`
	Name               string          `gorm:"type:varchar(255);not null" json:"product_name"`
	Category           FoodCategory    `gorm:"embedded" json:"food_category"`
	Price              decimal.Decimal `gorm:"type:numeric;not null" json:"food_price"`
	IsAvailable        bool            `gorm:"not null;default:true" json:"is_available"`
	ProductImg         []string        `gorm:"type:text[]" json:"product_img"`
	Description        string          `gorm:"type:text" json:"description"`
	Discount           Discount        `gorm:"embedded" json:"discount"`
	CreatedAt          time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
	AvgPreparationTime string          `gorm:"type:varchar(255);not null" json:"avg_preparation_time"`
	MaxPreparationTime string          `gorm:"type:varchar(255);not null" json:"max_preparation_time"`
	ShopId             uuid.UUID       `gorm:"type:uuid;not null;index" json:"shop_id"`
	Shop               Shop            `gorm:"foreignKey:ShopId;references:ShopId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"shop"`
	ProductType        string          `gorm:"type:varchar(50);not null" json:"product_type"`
	MaxQuantity        int             `gorm:"not null;default:100" json:"max_quantity"`
	MinQuantity        int             `gorm:"not null;default:1" json:"min_quantity"`
	AvailableQuantity  int             `gorm:"not null;default:100" json:"available_quantity"`
}
