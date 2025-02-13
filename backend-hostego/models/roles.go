package models

type Role struct {
	RoleId   string `gorm:"type:uuid;not null;unique;default:gen_random_uuid();" json:"role_id"`
	RoleName string `gorm:"type:varchaar(25);" json:"role_name"`

}
