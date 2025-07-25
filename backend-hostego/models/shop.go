package models

import "time"

type Shop struct {
	ShopId                 int          `gorm:"type:int;primaryKey;autoIncrement:true" json:"shop_id"`
	ShopName               string       `gorm:"type:varchar(255);" json:"shop_name"`
	ShopImg                string       `gorm:"type:varchar(255);" json:"shop_img"`
	Address                string       `gorm:"type:varchar(255);" json:"address"`
	PreparationTime        string       `gorm:"type:varchar(20);"  json:"preparation_time"`
	FoodCategory           FoodCategory `gorm:"embedded" json:"food_category"`
	ShopStatus             int          `gorm:"type:int;default:1;index;" json:"shop_status"`
	CreatedAt              time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	OwnerName              string       `gorm:"type:varchar(255);" json:"owner_name"`
	OwnerPhone             string       `gorm:"type:varchar(255);" json:"owner_phone"`
	OwnerEmail             string       `gorm:"type:varchar(255);" json:"owner_email"`
	Latitude               float64      `gorm:"type:double precision;" json:"latitude"`
	Longitude              float64      `gorm:"type:double precision;" json:"longitude"`
	ShopType               string       `gorm:"type:varchar(255);" json:"shop_type"`
	ShopDescription        string       `gorm:"type:varchar(255);" json:"shop_description"`
	ShopRating             float64      `gorm:"type:double precision;" json:"shop_rating"`
	FssaiLicenseNumber     string       `gorm:"type:varchar(255);" json:"fssai_license_number"`
	GSTINnumber            string       `gorm:"type:varchar(255);" json:"gstin_number"`
	FssaiLiscenseCopy      string       `gorm:"type:varchar(255);" json:"fssai_liscense_copy"`
	GSTINCopy              string       `gorm:"type:varchar(255);" json:"gstin_copy"`
	BankName               string       `gorm:"type:varchar(255);" json:"bank_name"`
	BankAccountNumber      string       `gorm:"type:varchar(255);" json:"bank_account_number"`
	BankIFSCCode           string       `gorm:"type:varchar(255);" json:"bank_ifsc_code"`
	BankAccountHolderName  string       `gorm:"type:varchar(255);" json:"bank_account_holder_name"`
	BankAccountType        string       `gorm:"type:varchar(255);" json:"bank_account_type"`
	Pancardcopy            string       `gorm:"type:varchar(255);" json:"pancard_copy"`
	ShopVerificationStatus string       `gorm:"type:varchar(255);" json:"shop_verification_status"`
	OwnerId                int          `gorm:"type:int;" json:"owner_id"`
	AverageRating          float64      `gorm:"type:double precision;default:0" json:"average_rating"`
	TotalRatings           int          `gorm:"type:int;default:0" json:"total_ratings"`
	OutletOpenTime         string       `gorm:"type:varchar(255);" json:"outlet_open_time"`
	OutletCloseTime        string       `gorm:"type:varchar(255);" json:"outlet_close_time"`
	// Owner                  User         `gorm:"foreignKey:OwnerId" json:"owner"`
}
