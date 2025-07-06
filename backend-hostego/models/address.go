package models

import "time"

type Address struct {
	AddressID    int       `gorm:"type:int;primaryKey;not null;autoIncrement:true" json:"address_id"`
	UserId       int       `gorm:"type:int;not null;index;" json:"user_id"`
	User         User      `gorm:"foreignKey:UserId;references:UserId;constraint:OnDelete:CASCADE" json:"user"`
	AddressType  string    `gorm:"type:varchar(255)" json:"address_type"`
	City         string    `gorm:"type:varchar(255)" json:"city"`
	PostalCode   string    `gorm:"type:varchar(255)" json:"postal_code"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	AddressLine1 string    `gorm:"type:varchar(255)" json:"address_line_1"`
}
