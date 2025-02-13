package models

import "time"

type UserRole struct {
	UserId    string    `gorm:"type:uuid;not null;index;" json:"user_id"`
	RoleId    string    `gorm:"type:uuid;not null;index;" json:"role_id"`
	User      User      `gorm:"type:uuid;foreignKey:user_id;" json:"user"`
	Role      Role      `gorm:"type:uuid;foreignKey:role_id;" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
