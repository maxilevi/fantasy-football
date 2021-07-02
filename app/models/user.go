package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email string
	PasswordHash []byte
	PermissionLevel int
}

func (u User) IsAdmin() bool {
	return u.PermissionLevel > 0
}