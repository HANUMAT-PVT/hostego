package dto

import (
	"backend-hostego/internal/app/hostego-service/models"
	"time"

	"github.com/google/uuid"
)

type ShopRequest struct {
	ShopId             uuid.UUID `json:"shop_id" binding:"required"`
	UserId             int       `json:"user_id"`
	ShopName           string    `json:"shop_name"`
	ShopImg            []string  `json:"shop_img"`
	Address            string    `json:"address"`
	AvgPreparationTime string    `json:"avg_preparation_time"`
	ShopEnabled        bool      `json:"shop_enabled"`
}

type UpdateShopRequest struct {
	ShopId             uuid.UUID           `json:"shop_id" binding:"required"`
	UserId             int                 `json:"user_id"`
	ShopName           string              `json:"shop_name"`
	ShopImg            []string            `json:"shop_img"`
	Address            string              `json:"address"`
	AvgPreparationTime string              `json:"avg_preparation_time"`
	ShopEnabled        bool                `json:"shop_enabled"`
	CreatedAt          *time.Time          `json:"created_at"`
	UpdatedAt          *time.Time          `json:"updated_at"`
	FoodCategory       models.FoodCategory `gorm:"embedded" json:"food_category"`
}

type CreateShopRequest struct {
	UserId             int                 `json:"user_id"`
	ShopName           string              `json:"shop_name"`
	ShopImg            []string            `json:"shop_img"`
	Address            string              `json:"address"`
	AvgPreparationTime string              `json:"avg_preparation_time"`
	ShopEnabled        bool                `json:"shop_enabled"`
	CreatedAt          *time.Time          `json:"created_at"`
	UpdatedAt          *time.Time          `json:"updated_at"`
	FoodCategory       models.FoodCategory `gorm:"embedded" json:"food_category"`
}
