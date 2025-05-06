package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserRequest struct {
	UserId              int        `json:"user_id" binding:"required"`
	FirstName           string     `json:"first_name"`
	LastName            string     `json:"last_name"`
	Email               string     `json:"email"`
	MobileNumber        string     `json:"mobile_number"`
	FirebaseOTPVerified int        `json:"firebase_otp_verified"`
	LastLoginTimestamp  time.Time  `json:"last_login_timestamp"`
	CreatedAt           *time.Time `json:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at"`
}

type UserAddressRequest struct {
	AddressId      uuid.UUID  `json:"address_id" binding:"required"`
	UserId         int        `json:"user_id" binding:"required"`
	AddressType    string     `json:"address_type"`
	City           string     `json:"city"`
	AddressLineOne string     `json:"address_line_one"`
	AddressLineTwo string     `json:"address_line_two"`
	PostalCode     string     `json:"postal_code"`
	State          string     `json:"state"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type UserAddressCreateRequest struct {
	UserId         int        `json:"user_id" binding:"required"`
	AddressType    string     `json:"address_type" binding:"required"`
	City           string     `json:"city" binding:"required"`
	AddressLineOne string     `json:"address_line_one" binding:"required"`
	AddressLineTwo string     `json:"address_line_two"`
	PostalCode     string     `json:"postal_code" binding:"required"`
	State          string     `json:"state" binding:"required"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}
