package models

import "time"

type UserRole struct {
	UserId    string    `gorm:"type:uuid;not null;index;" json:"user_id"`
	RoleId    string    `gorm:"type:uuid;not null;index;" json:"role_id"`
	User      User      `gorm:"type:foreignKey:User;" json:"user"`
	Role      Role      `gorm:"type:foreignKey:Role;" json:"role"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
