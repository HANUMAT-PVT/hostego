package dto

import (
	"backend-hostego/internal/app/hostego-service/models"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateProductRequestDto struct {
	Name               string              `json:"product_name"`
	Category           models.FoodCategory `json:"food_category"`
	Price              decimal.Decimal     `json:"food_price"`
	IsAvailable        bool                `json:"is_available"`
	ProductImg         []string            `json:"product_img"`
	Description        string              `json:"description"`
	Discount           models.Discount     `json:"discount"`
	CreatedAt          time.Time           `json:"created_at"`
	UpdatedAt          time.Time           `json:"updated_at"`
	AvgPreparationTime string              `json:"avg_preparation_time"`
	MaxPreparationTime string              `json:"max_preparation_time"`
	ShopId             uuid.UUID           `json:"shop_id"`
	ProductType        string              `json:"product_type"`
	MaxQuantity        int                 `json:"max_quantity"`
	MinQuantity        int                 `json:"min_quantity"`
	AvailableQuantity  int                 `json:"available_quantity"`
}
