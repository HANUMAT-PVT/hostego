package models

import "time"

type UserRole struct {
	UserRoleId int       `gorm:"type:int;not null;primaryKey;autoIncrement:true;unique;" json:"user_role_id"`
	UserId     string    `gorm:"type:string;not null;index;" json:"user_id"`
	RoleId     int       `gorm:"type:int;not null;index;" json:"role_id"`
	User       User      `gorm:"foreignKey:UserId;references:UserId" json:"user"`
	Role       Role      `gorm:"foreignKey:RoleId;references:RoleId" json:"role"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}
