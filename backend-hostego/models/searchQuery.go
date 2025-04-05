package models

import "time"

type SearchQuery struct {
	Id        int       `gorm:"type:int;primaryKey;autoIncrement;" json:"id"`
	Query     string    `gorm:"type:string;not null;" json:"query"`
	UserId    int       `gorm:"type:int;" json:"user_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
