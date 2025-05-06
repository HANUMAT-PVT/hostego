package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CartItemRequest struct {
	ProductId uuid.UUID       `json:"product_id"  binding:"required"`
	Quantity  int             `gorm:"not null" json:"quantity"  binding:"required"`
	UserId    int             `json:"user_id"`
	SubTotal  decimal.Decimal `json:"sub_total"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
}
