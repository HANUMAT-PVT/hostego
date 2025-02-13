package models

import "time"

type Order struct {
	OrderId   string    `gorm:"type:uuid;primaryKey;not null;unique;default:gen_random_uuid();" json:"order_id"`
	UserId    string    `gorm:"type:uuid;not null;unique;index;" json:"user_id"`
	User      User      `gorm:"type:uuid;foreignKey:User;" json:"user"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ShopId    string    `gorm:"type:uuid;not null;index;" json:"shop_id"`
	Shop      Shop      `gorm:"type:foreignKey:Shop;" json:"shop"`
}
