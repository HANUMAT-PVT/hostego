package models

import (
	"time"

	"github.com/google/uuid"
)

// User model

type User struct {
	UserId              int       `gorm:"primaryKey;autoIncrement" json:"user_id"`
	UserUuid            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()" json:"user_uuid"`
	UserCustodyWalletId int       `gorm:"type:int;" json:"user_custody_wallet_id"`
	Gender              string    `gorm:"type:varchar(255);" json:"gender""`
	FirstName           string    `gorm:"type:varchar(255);" json:"first_name"`
	LastName            string    `gorm:"type:varchar(255);" json:"last_name"`
	Email               string    `gorm:"type:text;unique;" json:"email"`
	MobileNumber        string    `gorm:"type:varchar(20);unique;not null;" json:"mobile_number"`
	FirebaseOTPVerified int       `gorm:"not null;default:0;" json:"firebase_otp_verified"`
	LastLoginTimestamp  time.Time `gorm:"autoUpdateTime" json:"last_login_timestamp"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DateOfBirth         string    `gorm:"type:string;" json:"date_of_birth"`
	IsVerifiedEmail     bool      `gorm:"not null;default:false;" json:"is_verified_email"`
	RefferalCode        string    `gorm:"type:varchar(255);" json:"refferal_code"`
	AdvertisementId     int       `gorm:"type:int;" json:"advertisement_id"`
	IsVerifiedMobile    bool      `gorm:"not null;default:false;" json:"is_verified_mobile"`
	IsUserBlocked       bool      `gorm:"not null;default:false;" json:"is_user_blocked"`
	IsWalletDisabled    bool      `gorm:"not null;default:false;" json:"is_wallet_disabled"`
	IsDeliveryPerson    bool      `gorm:"not null;default:false;" json:"is_delivery_person"`
	Rating              float32   `gorm:"default:0.0;" json:"rating"`
}

// Before create hook to generate uuid
// func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
// 	u.UserId = uuid.New().String()
// 	return
// }

// func (u *User) AfterCreate(tx *gorm.DB) (err error) {
// 	wallet := Wallet{
// 		WalletId: uuid.New().String(),
// 		UserId:   u.UserId,
// 		Balance:  0.0,
// 	}
// 	if err := tx.Create(&wallet).Error; err != nil {
// 		return err
// 	}

// 	return nil

// }
