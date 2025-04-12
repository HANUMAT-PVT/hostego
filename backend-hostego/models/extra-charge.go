package models

type ExtrCharge struct {
	AddressID int `gorm:"type:int;primaryKey;not null;autoIncrement:true" json:"address_id"`
}
