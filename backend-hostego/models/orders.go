package models

import (
	"time"

	"gorm.io/datatypes"
)

type OrderStatusType string

const (
	PendingOrderStatus   = "pending"
	PlacedOrderStatus    = "placed"
	AssignedOrderStatus  = "assigned"
	PackedOrderStatus    = "packed"
	ReachedOrderStatus   = "reached"
	PickedOrderStatus    = "picked"
	OnTheWayOrderStatus  = "on_the_way"
	DeliveredOrderStatus = "delivered"
	CookingOrderStatus   = "cooking"
	CanceledOrderStatus  = "cancelled"
)

type Order struct {
	OrderId              string             `gorm:"type:string;primaryKey;not null;unique;default:gen_random_uuid();" json:"order_id"`
	UserId               string             `gorm:"type:text;not null;index;" json:"user_id"`
	User                 User               `gorm:"foreignKey:UserId;references:UserId" json:"user"`
	CreatedAt            time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	OrderItems           datatypes.JSON     `gorm:"type:jsonb" json:"order_items"`
	PlatformFee          float64            `gorm:"type:double precision;not null;" json:"platform_fee"`
	ShippingFee          float64            `gorm:"type:double precision;not null;" json:"shipping_fee"`
	FinalOrderValue      float64            `gorm:"type:double precision;not null;" json:"final_order_value"`
	DeliveryPartnerFee   float64            `gorm:"type:double precision;not null;" json:"delivery_partner_fee"`
	PaymentTransactionId string             `gorm:"type:string" json:"payment_transaction_id"`
	PaymentTransaction   PaymentTransaction `gorm:"foreignKey:PaymentTransactionId;references:PaymentTransactionId" json:"payment_transaction"`
	OrderStatus          OrderStatusType    `gorm:"type:varchar(20);default:pending" json:"order_status"`
	DeliveryPartner      datatypes.JSON     `gorm:"type:jsonb" json:"delivery_partner"`
	DeliveryPartnerId    string             `gorm:"type:string" json:"delivery_partner_id"`
	AddressID            string             `gorm:"type:string;not null" json:"address_id"`
	Address              Address            `gorm:"foreignKey:AddressID;references:AddressID" json:"address"`
	DeliveredAt          time.Time          `json:"delivered_at"`
	RefundedAt           time.Time          `json:"refunded_at"`
	Refunded             bool               `gorm:"type:boolean;default:false" json:"refunded"`
	RefundInitiator      string             `gorm:"type:string" json:"refund_initiator"`
	FreeDelivery         bool               `gorm:"type:boolean;default:false" json:"free_delivery"`
	CookingRequests      string             `gorm:"type:string" json:"cooking_requests"`
}
