package models

import (
	"gorm.io/gorm"
)

// User DB model
type User struct {
	gorm.Model
	Email           string
	PasswordHash    []byte
	PermissionLevel int
}

// Returns a bool that represents if the user has admin privileges.
func (u User) IsAdmin() bool {
	return u.PermissionLevel > 0
}

type ShowUser struct {
	ID    uint     `json:"id"`
	Email string   `json:"email"`
	Team  ShowTeam `json:"team"`
} //@name ShowUser

type UpdateUser struct {
	Email string `json:"email"`
} //@name UpdateUser

type CreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
} //@name CreateUser
