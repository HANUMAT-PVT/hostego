package models;

import (
	"time"
)

type MessMenu struct {
	ID int `gorm:"type:int;primary_key;not null;autoIncrement:true;" json:"id"`
	Date string `gorm:"type:date;not null;" json:"date"`
	Menu string `gorm:"type:text;not null;" json:"menu"`
	CreatedAt time.Time `gorm:"autoCreateTime;" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;" json:"updated_at"`
	
}
