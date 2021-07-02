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
	user := models.User{
		Email: email,
		PasswordHash: hash,
		PermissionLevel: permission,
	}
	u.Db.Create(&user)
	team, players := models.RandomTeam()
	team.OwnerID = user.ID
	u.Db.Create(&team)
	for i := range players {
		players[i].TeamID = team.ID
		u.Db.Create(&players[i])
	}
}

func (u RepositorySQL) GetUser(email string, user *models.User) error {
	res := u.Db.Where(models.User{Email: email}).First(&user)
	return res.Error
}

type RepositoryMemory struct {
	Users []models.User
	Teams map[string]models.Team
	Players map[string]models.Player
}

func CreateRepositoryMemory() *RepositoryMemory {
	return &RepositoryMemory{
		Users: make([]models.User, 0),
		Teams: map[string]models.Team{},
		Players: map[string]models.Player{},
	}
}

func (u *RepositoryMemory) CreateUser(email string, hash []byte, permission int) {
	u.Users = append(u.Users, models.User{
		Email: email,
		PasswordHash: hash,
		PermissionLevel: permission,
	})
	team, players := models.RandomTeam()
	u.Teams[email] = team
	for i := range players {
		u.Players[email] = players[i]
	}
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