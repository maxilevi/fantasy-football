package repos

import (
	"fmt"
	"gorm.io/gorm"
)
import "../models"

type Repository interface {
	CreateUser(email string, hash []byte, permission int)
	GetUser(email string, user *models.User) error
	GetTeam(id uint) (models.Team, error)
	GetPlayer(playerId uint) (models.Player, error)
	GetPlayers(teamId uint) []models.Player
	UpdateTeam(team models.Team)
	UpdatePlayer(player models.Player)
	GetUserTeam(user models.User) (models.Team, error)
}

type RepositorySQL struct {
	Db *gorm.DB
}

func (u RepositorySQL) CreateUser(email string, hash []byte, permission int) {
	user := models.User{
		Email:           email,
		PasswordHash:    hash,
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

func (u RepositorySQL) GetTeam(id uint) (models.Team, error) {
	var team models.Team
	res := u.Db.First(&team, id)
	return team, res.Error
}

func (u RepositorySQL) GetPlayers(teamId uint) []models.Player {
	var players []models.Player
	u.Db.Where(&models.Player{TeamID: teamId}).Find(&players)
	return players
}

func (u RepositorySQL) GetPlayer(playerId uint) (models.Player, error) {
	var player models.Player
	res := u.Db.Find(&player, playerId)
	return player, res.Error
}

func (u RepositorySQL) UpdateTeam(team models.Team) {
	u.Db.Save(team)
}

func (u RepositorySQL) UpdatePlayer(player models.Player) {
	u.Db.Save(player)
}

func (u RepositorySQL) GetUserTeam(user models.User) (models.Team, error) {
	var team models.Team
	res := u.Db.Where(&models.Team{OwnerID: user.ID}).Find(&team)
	return team, res.Error
}

type RepositoryMemory struct {
	Users   []models.User
	Teams   map[string]models.Team
	Players map[string]models.Player
}

func CreateRepositoryMemory() *RepositoryMemory {
	return &RepositoryMemory{
		Users:   make([]models.User, 0),
		Teams:   make(map[string]models.Team),
		Players: make(map[string]models.Player),
	}
}

func (u *RepositoryMemory) CreateUser(email string, hash []byte, permission int) {
	u.Users = append(u.Users, models.User{
		Email:           email,
		PasswordHash:    hash,
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

func (u *RepositoryMemory) GetTeam(id uint) (models.Team, error) {
	panic("implement me")
}

func (u *RepositoryMemory) GetPlayer(playerId uint) (models.Player, error) {
	panic("implement me")
}

func (u *RepositoryMemory) GetPlayers(teamId uint) []models.Player {
	panic("implement me")
}

func (u *RepositoryMemory) UpdateTeam(team models.Team) {
	panic("implement me")
}

func (u *RepositoryMemory) UpdatePlayer(player models.Player) {
	panic("implement me")
}

func (u *RepositoryMemory) GetUserTeam(user models.User) (models.Team, error) {
	panic("implement me")
}
