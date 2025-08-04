package models

import "time"

type Notification struct {
	Id                   int       `gorm:"type:int;primaryKey;autoIncrement:true" json:"id"`
	Title                string    `gorm:"type:varchar(255);" json:"title"`
	Body                 string    `gorm:"type:varchar(255);" json:"body"`
	Link                 string    `gorm:"type:varchar(255);" json:"link"`
	NotificationImageUrl string    `gorm:"type:varchar(255);" json:"notification_image_url"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
