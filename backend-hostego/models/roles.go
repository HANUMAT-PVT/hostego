package models

type Role struct {
	RoleId   int    `gorm:"type:int;primaryKey;not null;unique;autoIncrement;" json:"role_id"`
	RoleName string `gorm:"type:varchar(25);" json:"role_name"`
}
