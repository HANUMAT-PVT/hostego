package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CartItem struct {
	CartItemId  uuid.UUID       `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"cart_item_id"`
	ProductId   uuid.UUID       `gorm:"type:uuid;not null" json:"product_id"`
	Quantity    int             `gorm:"not null" json:"quantity"`
	SubTotal    decimal.Decimal `gorm:"type:numeric;not null" json:"sub_total"`
	UserId      int64           `gorm:"not null;index" json:"user_id"`
	ProductItem Product         `gorm:"foreignKey:ProductId;references:ProductId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product_item"`
	Status      string          `gorm:"type:varchar(50);default:'pending'" json:"status"`
	CreatedAt   time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}
