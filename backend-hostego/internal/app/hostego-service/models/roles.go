package models

import (
	
)

type Role struct {
	RoleId   string `gorm:"type:uuid;primaryKey;not null;unique;default:gen_random_uuid();" json:"role_id"`
	RoleName string `gorm:"type:varchar(25);" json:"role_name"`

}



