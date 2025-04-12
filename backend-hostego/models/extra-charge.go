package models

import "time"

type ExtrCharge struct {
	Id         int       `gorm:"type:int;primaryKey;not null;autoIncrement:true" json:"id"`
	ChargeType string    `gorm:"type:string;" json:"charge_type"`
	Amount     float64   `gorm:"type:float;" json:"charge_amount"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Status     string    `gorm:"type:string;" json:"status_type"`
}
