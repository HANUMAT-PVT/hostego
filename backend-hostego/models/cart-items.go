package models

type CartItem struct {
	CartItemId    string  `gorm:"type:uuid;not null;unique;primaryKey;default:gen_random_uuid()" json:"cart_item_id"`
	ProductItemID string  `gorm:"type:varchar(50);not null;index" json:"product_item"`
	Quantity      float64 `gorm:"type:double precision;default:1;not null;" json:"quantity"`
	SubTotal      float64 `gorm:"type:double precision;not null;" json:"sub_total"`

	// Foreign Key Relation
	ProductItem   Product `gorm:"foreignKey:product_id;references:Product" json:"product_item"` 
}
