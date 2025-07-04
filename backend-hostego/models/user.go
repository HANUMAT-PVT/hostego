package models

import (
	"time"

	"gorm.io/gorm"
)

// User model
type User struct {
	UserId              int       `gorm:"type:int;primaryKey;unique;not null;autoIncrement:true" json:"user_id"`
	FirstName           string    `gorm:"type:varchar(255);" json:"first_name"`
	LastName            string    `gorm:"type:varchar(255);" json:"last_name"`
	Email               string    `gorm:"unique;not null;" json:"email"`
	MobileNumber        string    `gorm:"type:varchar(20);unique;not null;" json:"mobile_number"`
	FirebaseOTPVerified int       `gorm:"not null;default:0;" json:"firebase_otp_verified"`
	LastLoginTimestamp  time.Time `gorm:"autoUpdateTime" json:"last_login_timestamp"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	FCMToken            string    `gorm:"type:varchar(255);" json:"fcm_token"`
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	wallet := Wallet{
		UserId:  u.UserId,
		Balance: 0.0,
	}
	if err := tx.Create(&wallet).Error; err != nil {
		return err
	}

	return nil

}
