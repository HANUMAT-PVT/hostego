package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type OrderStatusType string

const (
	PendingOrderStatus   = "PENDING"
	PackedOrderStatus    = "PACKED"
	CookingOrderStatus   = "COOKING"
	DeliveredOrderStatus = "DELIVERED"
	CanceledOrderStatus  = "CANCELED"
	PickedOrderStatus    = "PICKED"
)

type Order struct {
	OrderId              string             `gorm:"type:uuid;primaryKey;not null;unique;default:gen_random_uuid();" json:"order_id"`
	UserId               string             `gorm:"type:text;not null;unique;index;" json:"user_id"`
	User                 User               `gorm:"foreignKey:UserId;references:UserId" json:"user"`
	CreatedAt            time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	OrderItems           datatypes.JSON     `gorm:"type:jsonb" json:"order_items"`
	PlatformFee          float64            `gorm:"type:double precision;not null;" json:"platform_fee"`
	ShippingFee          float64            `gorm:"type:double precision;not null;" json:"shipping_fee"`
	FinalOrderValue      float64            `gorm:"type:double precision;not null;" json:"final_order_value"`
	DeliveryPartnerFee   float64            `gorm:"type:double precision;not null;" json:"delivery_partner_fee"`
	PaymentTransactionId uuid.UUID          `gorm:"type:uuid" json:"payment_transaction_id"`
	PaymentTransaction   PaymentTransaction `gorm:"foreignKey:PaymentTransactionId;references:PaymentTransactionId" json:"payment_transaction"`
	OrderStatus          OrderStatusType    `gorm:"type:varchar(20);default:pending" json:"order_status"`
	DeliveryPartnerId    datatypes.JSON     `gorm:"type:jsonb" json:"delivery_partner_id"`

	DeliveredAt time.Time `json:"delivered_at"`
}
