package models

import (
	"time"

	"github.com/google/uuid"
)

type Documents struct {
	AadhaarFrontImg string `json:"aadhaar_front_img"` // Link to Aadhaar front image
	AadhaarBackImg  string `json:"aadhaar_back_img"`  // Link to Aadhaar back image
	UPI_ID          string `json:"upi_id"`            // UPI ID (if it's a link)
	BankDetailsImg  string `json:"bank_details_img"`  // Link to bank details image
}

type DeliveryPartner struct {
	DeliveryPartnerID  uuid.UUID `gorm:"type:uuid;primaryKey;not null;default:gen_random_uuid();" json:"delivery_partner_id"`
	UserId             string    `gorm:"type:string;not null;unique;" json:"user_id"`
	User               User      `gorm:"foreignKey:UserId;references:UserId;" json:"user"`
	AvailabilityStatus int       `gorm:"type:int;not null;default:0;" json:"availability_status"`
	Address            string    `gorm:"type:varchar(255);" json:"address"`
	AccountStatus      int       `gorm:"type:int;default:0;" json:"account_status"`
	Documents          Documents `gorm:"embedded" json:"documents"`
	VerificationStatus int       `gorm:"type:int;default:0;" json:"verification_status"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
