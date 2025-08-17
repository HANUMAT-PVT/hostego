package models

import (
	"time"

	"gorm.io/datatypes"
)

type OrderStatusType string

type OrderType string

const (
	DeliveryOrderType = "delivery"
	TakeawayOrderType = "takeaway"
)

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
	ReachedDoorStatus    = "reached_door"
	CanceledOrderStatus  = "cancelled"
	PreparingOrderStatus = "preparing"
	ReadyOrderStatus     = "ready"
	AcceptedByRestaurant = "accepted_by_restaurant"
	RejectedByRestaurant = "rejected_by_restaurant"
)

type Order struct {
	OrderId                 int                `gorm:"type:int;primaryKey;not null;unique;autoIncrement:true" json:"order_id"`
	UserId                  int                `gorm:"type:int;not null;index;" json:"user_id"`
	User                    User               `gorm:"foreignKey:UserId;references:UserId" json:"user"`
	CreatedAt               time.Time          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt               time.Time          `gorm:"autoUpdateTime" json:"updated_at"`
	OrderItems              datatypes.JSON     `gorm:"type:jsonb" json:"order_items"`
	PlatformFee             float64            `gorm:"type:double precision;not null;" json:"platform_fee"`
	ShippingFee             float64            `gorm:"type:double precision;not null;" json:"shipping_fee"`
	FinalOrderValue         float64            `gorm:"type:double precision;not null;" json:"final_order_value"`
	DeliveryPartnerFee      float64            `gorm:"type:double precision;not null;" json:"delivery_partner_fee"`
	PaymentTransactionId    int                `gorm:"type:int" json:"payment_transaction_id"`
	PaymentTransaction      PaymentTransaction `gorm:"foreignKey:PaymentTransactionId;references:PaymentTransactionId" json:"payment_transaction"`
	OrderStatus             OrderStatusType    `gorm:"type:varchar(20);default:pending" json:"order_status"`
	DeliveryPartner         datatypes.JSON     `gorm:"type:jsonb" json:"delivery_partner"`
	DeliveryPartnerId       int                `gorm:"type:int" json:"delivery_partner_id"`
	AddressID               int                `gorm:"type:int;not null" json:"address_id"`
	Address                 Address            `gorm:"foreignKey:AddressID;references:AddressID" json:"address"`
	DeliveredAt             time.Time          `json:"delivered_at"`
	RefundedAt              time.Time          `json:"refunded_at"`
	Refunded                bool               `gorm:"type:boolean;default:false" json:"refunded"`
	RefundInitiator         int                `gorm:"type:int" json:"refund_initiator"`
	FreeDelivery            bool               `gorm:"type:boolean;default:false" json:"free_delivery"`
	CookingRequests         string             `gorm:"type:string" json:"cooking_requests"`
	ShopId                  int                `gorm:"type:int" json:"shop_id"`
	Shop                    Shop               `gorm:"foreignKey:ShopId;references:ShopId" json:"shop"`
	RestaurantRespondedAt   time.Time          `gorm:"type:timestamp" json:"restaurant_responded_at"`
	ExpectedReadyAt         time.Time          `gorm:"type:timestamp" json:"expected_ready_at"`
	ActualReadyAt           time.Time          `gorm:"type:timestamp" json:"actual_ready_at"`
	RejectedByRestaurant    bool               `gorm:"type:boolean;default:false" json:"rejected_by_restaurant"`
	IsAcceptedByRestaurant  bool               `gorm:"type:boolean;default:false" json:"is_accepted_by_restaurant"`
	IsRejectedByRestaurant  bool               `gorm:"type:boolean;default:false" json:"is_rejected_by_restaurant"`
	RestaurantPayableAmount float64            `gorm:"type:double precision;default:0" json:"restaurant_payable_amount"`
	RestaurantPaidAt        *time.Time         `gorm:"type:timestamp;default:null" json:"restaurant_paid_at"`
	RestaurantPayoutID      string             `gorm:"type:varchar(255)" json:"restaurant_payout_id"`
	OrderType               OrderType          `gorm:"type:varchar(20);default:delivery" json:"order_type"`
}
