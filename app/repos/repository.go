package repos

import (
	"fmt"
	"gorm.io/gorm"
)
import "../models"

type Repository interface {
	CreateUser(email string, hash []byte, permission int)
	GetUser(email string, user *models.User) error
}

type RepositorySQL struct {
	Db *gorm.DB
}

func (u RepositorySQL) CreateUser(email string, hash []byte, permission int) {
	u.Db.Create(&models.User{
		Email: email,
		PasswordHash: hash,
		PermissionLevel: permission,
	})
}

func (u RepositorySQL) GetUser(email string, user *models.User) error {
	res := u.Db.Where(models.User{Email: email}).First(&user)
	return res.Error
}

type RepositoryMemory struct {
	Users []models.User
}

func (u *RepositoryMemory) CreateUser(email string, hash []byte, permission int) {
	u.Users = append(u.Users, models.User{
		Email: email,
		PasswordHash: hash,
		PermissionLevel: permission,
	})
}

func (u *RepositoryMemory) GetUser(email string, user *models.User) error {
	for i := range u.Users {
		if u.Users[i].Email == email {
			*user = u.Users[i]
			return nil
		}
	}
	return fmt.Errorf("user not found")
}