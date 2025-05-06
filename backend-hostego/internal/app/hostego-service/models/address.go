package models

import (
	"time"

	"github.com/google/uuid"
)

type UserAddress struct {
	AddressId      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"address_id"`
	UserId         int       `gorm:"not null;index" json:"user_id"`
	User           User      `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE" json:"-"`
	AddressType    string    `gorm:"type:varchar(255)" json:"address_type"`
	City           string    `gorm:"type:varchar(255)" json:"city"`
	AddressLineOne string    `gorm:"type:varchar(255)" json:"address_line_one"`
	AddressLineTwo string    `gorm:"type:varchar(255)" json:"address_line_two"`
	PostalCode     string    `gorm:"type:varchar(255)" json:"postal_code"`
	State          string    `gorm:"type:varchar(255)" json:"state"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
