package models

import "time"

type Notification struct {
	Id        int       `gorm:"type:int;primaryKey;autoIncrement:true" json:"id"`
	Title     string    `gorm:"type:varchar(255);" json:"title"`
	Message   string    `gorm:"type:varchar(255);" json:"message"`
	UserId    int       `gorm:"type:int;not null;" json:"user_id"`
	IsRead    bool      `gorm:"type:boolean;default:false" json:"is_read"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
}
