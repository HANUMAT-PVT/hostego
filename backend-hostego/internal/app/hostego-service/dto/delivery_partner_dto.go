package dto

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryPartnerRequest struct {
	DeliveryPartnerId uuid.UUID `json:"delivery_partner_id"`
	UserId            int       `json:"user_id"`
}
type DeliveryPartnerRequestDto struct {
	UserId          int        `json:"user_id"`
	IsAvailable     bool       `json:"is_available"`
	Address         string     `json:"address"`
	AccountVerified bool       `json:"account_verified"`
	VehicleType     string     `json:"vehicle_type"`
	VehicleNumber   string     `json:"vehicle_number"`
	PartnerImgUrl   string     `json:"partner_img_url"`
	IsEnabled       bool       `json:"is_enabled"`
	Rating          float32    `json:"rating"`
	CreatedAt       *time.Time `json:"created_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
}
