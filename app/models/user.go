package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email           string
	PasswordHash    []byte
	PermissionLevel int
}

func (u User) IsAdmin() bool {
	return u.PermissionLevel > 0
}

type ShowUser struct {
	ID    uint     `json:"id"`
	Email string   `json:"email"`
	Team  ShowTeam `json:"team"`
}

type UpdateUser struct {
	Email string `json:"email"`
	Team  uint   `json:"team"`
}

type CreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
