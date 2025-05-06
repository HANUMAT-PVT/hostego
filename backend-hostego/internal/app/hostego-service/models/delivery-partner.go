package models

import (
	"time"

	"github.com/google/uuid"
)

// Documents struct stores all required documents for a Delivery Partner
type Documents struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"id"`
	DeliveryPartnerId uuid.UUID `gorm:"type:uuid;not null;index;" json:"delivery_partner_id"`
	AadhaarFrontImg   string    `gorm:"type:varchar(255);" json:"aadhaar_front_img"`
	AadhaarBackImg    string    `gorm:"type:varchar(255);" json:"aadhaar_back_img"`
	UPI_ID            string    `gorm:"type:varchar(255);" json:"upi_id"`
	BankDetailsImg    string    `gorm:"type:varchar(255);" json:"bank_details_img"`
	LicenseImg        string    `gorm:"type:varchar(255);" json:"license_img"`
	IsVerifiedAll     bool      `gorm:"default:false;" json:"is_verified_all"`
}

type DeliveryPartner struct {
	DeliveryPartnerId uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid();" json:"delivery_partner_id"`
	UserId            int       `gorm:"not null;index;" json:"user_id"`
	User              User      `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`
	IsAvailable       bool      `gorm:"default:false;" json:"is_available"`
	Address           string    `gorm:"type:varchar(255);" json:"address"`
	AccountVerified   bool      `gorm:"default:false;" json:"account_verified"`
	VehicleType       string    `gorm:"type:varchar(100);" json:"vehicle_type"`
	VehicleNumber     string    `gorm:"type:varchar(50);unique;not null;" json:"vehicle_number"`
	PartnerImgUrl     string    `gorm:"type:varchar(255);" json:"partner_img_url"`
	IsEnabled         bool      `gorm:"default:true;" json:"is_enabled"`
	Rating            float32   `gorm:"default:0.0;" json:"rating"`
	Documents         Documents `gorm:"foreignKey:DeliveryPartnerId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"documents"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
