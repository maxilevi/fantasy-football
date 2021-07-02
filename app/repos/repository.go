package repos

import "gorm.io/gorm"
import "../models"

type Repository interface {
	Create(email string, hash []byte, permission int)
	GetUser(options *models.User, user *models.User) error
}

type RepositorySQL struct {
	Db *gorm.DB
}

func (u RepositorySQL) Create(email string, hash []byte, permission int) {
	u.Db.Create(&models.User{
		Email: email,
		PasswordHash: hash,
		PermissionLevel: permission,
	})
}

func (u RepositorySQL) GetUser(options *models.User, user *models.User) error {
	res := u.Db.Where(options).First(&user)
	return res.Error
}