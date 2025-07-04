package models

type ProductCategory struct {
	CategoryId   int    `gorm:"type:int;primaryKey;autoIncrement:true" json:"category_id"`
	CategoryName string `gorm:"type:varchar(255);not null" json:"category_name"`
	ShopId       int    `gorm:"type:int;not null;index" json:"shop_id"`
	Shop         Shop   `gorm:"foreignKey:ShopId;references:ShopId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"shop"`
}
