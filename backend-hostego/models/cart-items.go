package models

import "github.com/google/uuid"

type CartItem struct {
	CartItemId string    `gorm:"type:uuid;not null;unique;primaryKey;default:gen_random_uuid()" json:"cart_item_id"`
	ProductId  uuid.UUID `gorm:"type:uuid;not null;index" json:"product_id"`
	Quantity   int       `gorm:"type:int;default:1;not null;" json:"quantity"`
	SubTotal   float64   `gorm:"type:double precision;not null;" json:"sub_total"`
	UserId     string    `gorm:"type:text;not null;index;" json:"user_id"`

	// Foreign Key Relation
	ProductItem Product `gorm:"foreignKey:ProductId;references:ProductId" json:"product_item"`
}
