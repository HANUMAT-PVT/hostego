package models

import (
	"time"

	"github.com/google/uuid"
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
	OrderId            string             `gorm:"type:uuid;primaryKey;not null;unique;default:gen_random_uuid();" json:"order_id"`
	UserId             string             `gorm:"type:uuid;not null;unique;index;" json:"user_id"`
	User               User               `gorm:"type:uuid;foreignKey:User;" json:"user"`
	CreatedAt          time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	ShopId             string             `gorm:"type:uuid;not null;index;" json:"shop_id"`
	Shop               Shop               `gorm:"type:foreignKey:ShopId;references:ShopId" json:"shop"`
	OrderItems         []CartItem         `gorm:"type:jsonb" json:"order_items"`
	PlatformFee        float64            `gorm:"type:double precision;default:1.0:not null;" json:"platform_fee"`
	ShippingFee        float64            `gorm:"type:double precision;not null;" json:"shipping_fee"`
	FinalOrderValue    float64            `gorm:"type:double precision;not null;" json:"final_order_value"`
	DeliveryPartnerFee float64            `gorm:"type:double precision:not null;" json:"delivery_partner_fee"`
	PaymentId          uuid.UUID          `gorm:"type:uuid" json:"payment_transaction_id"`
	PaymentTransaction PaymentTransaction `gorm:"foreignKey:PaymentId;references:PaymentId" json:"payment_transaction"`
	OrderStatus        OrderStatusType    `gorm:"type:varchar(20);" json:"order_status"`
	DeliveryPartnerId  uuid.UUID          `gorm:"type:uuid" json:"delivery_partner_id"`
	DeliveryPartner    DeliveryPartner    `gorm:"foreignKey:DeliveryPartnerId;references:DeliveryPartnerId;" json:"delivery_partner"`

	DeliveredAt time.Time `json:"delivered_at"`
}
