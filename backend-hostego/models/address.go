package models

import "time"

type Address struct {
	AddressID   string    `gorm:"type:uuid;unique;not null;primaryKey;default:gen_random_uuid()" json:"address_id"`
	UserId      string    `gorm:"type:uuid;not null; index;" json:"user_id"`
	User        User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"user"`
	AddressType string    `gorm:"type:varchar(255)" json:"address_type"`
	City        string    `gorm:"type:varchar(255)" json:"city"`
	PostalCode  string    `gorm:"type:varchar(255)" json:"postal_code"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
