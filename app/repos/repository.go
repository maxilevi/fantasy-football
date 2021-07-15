package repos

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"time"
)
import "../models"

// Repository pattern to handle abstraction of the data source
type Repository interface {
	CreateUser(email string, hash []byte, permission int) (models.User, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserById(id uint) (models.User, error)
	GetTeam(id uint) (models.Team, error)
	GetPlayer(playerId uint) (models.Player, error)
	GetPlayers(teamId uint) []models.Player
	GetUserTeam(user models.User) (models.Team, error)
	Create(model interface{}) error
	Update(model interface{}) error
	Delete(model interface{}) error
	GetTransfers() []models.Transfer
	GetTransfer(id uint) (models.Transfer, error)
	RunInTransaction(code func() error) error
	DeleteTeam(team *models.Team) error
	DeletePlayer(player *models.Player) error
	GetTransferWithPlayer(player *models.Player) (models.Transfer, error)
}

// Create an user on a given repository
func doCreateUser(u Repository, email string, hash []byte, permission int) (models.User, error) {
	user := models.User{
		Email:           email,
		PasswordHash:    hash,
		PermissionLevel: permission,
	}
	return user, u.RunInTransaction(func() error {
		err := u.Create(&user)
		if err != nil {
			return err
		}

		team, players := models.RandomTeam()
		team.UserID = user.ID
		err = u.Create(&team)
		if err != nil {
			return err
		}

		for i := range players {
			players[i].TeamID = team.ID
			err = u.Create(&players[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Delete a team on a given repository
func doDeleteTeam(u Repository, team *models.Team) error {
	return u.RunInTransaction(func() error {
		players := u.GetPlayers(team.ID)
		for _, p := range players {
			err := u.DeletePlayer(&p)
			if err != nil {
				return err
			}
		}
		return u.Delete(team)
	})
}

// Delete a player on a given repository
func doDeletePlayer(u Repository, player *models.Player) error {
	return u.RunInTransaction(func() error {
		transfer, err := u.GetTransferWithPlayer(player)
		if err == nil {
			// Transfer exists, delete it
			if err := u.Delete(&transfer); err != nil {
				return err
			}
		}
		return u.Delete(player)
	})
}

// Implementation of the repository interface using a DB connection
type RepositorySQL struct {
	Db *gorm.DB
}

// Create a new user
func (u RepositorySQL) CreateUser(email string, hash []byte, permission int) (models.User, error) {
	return doCreateUser(u, email, hash, permission)
}

// Get an user by email
func (u RepositorySQL) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	res := u.Db.Where(models.User{Email: email}).First(&user)
	if res.Error == nil && user.CreatedAt == (time.Time{}) {
		return user, fmt.Errorf("record not found")
	}
	return user, res.Error
}

// Get an user by id
func (u RepositorySQL) GetUserById(id uint) (models.User, error) {
	var user models.User
	res := u.Db.Preload(clause.Associations).First(&user, id)
	if res.Error == nil && user.CreatedAt == (time.Time{}) {
		return user, fmt.Errorf("record not found")
	}
	return user, res.Error
}

// Get a team by id
func (u RepositorySQL) GetTeam(id uint) (models.Team, error) {
	var team models.Team
	res := u.Db.Preload(clause.Associations).First(&team, id)
	if res.Error == nil && team.CreatedAt == (time.Time{}) {
		return team, fmt.Errorf("record not found")
	}
	return team, res.Error
}

// Get a players from a specific team
func (u RepositorySQL) GetPlayers(teamId uint) []models.Player {
	var players []models.Player
	u.Db.Preload(clause.Associations).Where(&models.Player{TeamID: teamId}).Find(&players)
	return players
}

// Get a player by id
func (u RepositorySQL) GetPlayer(playerId uint) (models.Player, error) {
	var player models.Player
	res := u.Db.Preload(clause.Associations).Find(&player, playerId)
	if res.Error == nil && player.CreatedAt == (time.Time{}) {
		return player, fmt.Errorf("record not found")
	}
	return player, res.Error
}

// Create a new record given a model
func (u RepositorySQL) Create(model interface{}) error {
	res := u.Db.Save(model)
	return res.Error
}

// Update a new record given a model
func (u RepositorySQL) Update(model interface{}) error {
	res := u.Db.Save(model)
	fmt.Println(res)
	return res.Error
}

// Delete a new record given a model
func (u RepositorySQL) Delete(model interface{}) error {
	res := u.Db.Unscoped().Delete(model)
	return res.Error
}

// Get an user's attached team
func (u RepositorySQL) GetUserTeam(user models.User) (models.Team, error) {
	var team models.Team
	res := u.Db.Preload(clause.Associations).Where(&models.Team{UserID: user.ID}).Find(&team)
	if res.Error == nil && team.CreatedAt == (time.Time{}) {
		return team, fmt.Errorf("record not found")
	}
	return team, res.Error
}

// Get all existing transfers
func (u RepositorySQL) GetTransfers() []models.Transfer {
	var transfers []models.Transfer
	u.Db.Preload("Player.Team").Where("1 = 1").Find(&transfers)
	return transfers
}

// Get a transfer by id
func (u RepositorySQL) GetTransfer(id uint) (models.Transfer, error) {
	var transfer models.Transfer
	res := u.Db.Preload("Player.Team").Find(&transfer, id)
	if res.Error == nil && transfer.CreatedAt == (time.Time{}) {
		return transfer, fmt.Errorf("record not found")
	}
	return transfer, res.Error
}

// Run the function inside a transaction and rollback in case of error
func (u RepositorySQL) RunInTransaction(code func() error) error {
	tx := u.Db.Begin()
	err := code()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

// Delete a given team
func (u RepositorySQL) DeleteTeam(team *models.Team) error {
	return doDeleteTeam(u, team)
}

// Delete a given player
func (u RepositorySQL) DeletePlayer(player *models.Player) error {
	return doDeletePlayer(u, player)
}

// Get a transfer with a player
func (u RepositorySQL) GetTransferWithPlayer(player *models.Player) (models.Transfer, error) {
	var transfer models.Transfer
	res := u.Db.Preload(clause.Associations).Where(&models.Transfer{PlayerID: player.ID}).Find(&transfer)
	if res.Error == nil && transfer.CreatedAt == (time.Time{}) {
		return transfer, fmt.Errorf("record not found")
	}
	return transfer, res.Error
}

// Repository implementation with models on memory
type RepositoryMemory struct {
	Models []interface{}
}

// Create a new memory repository
func CreateRepositoryMemory() *RepositoryMemory {
	return &RepositoryMemory{
		Models: make([]interface{}, 0),
	}
}

// Create a user
func (u *RepositoryMemory) CreateUser(email string, hash []byte, permission int) (models.User, error) {
	return doCreateUser(u, email, hash, permission)
}

// Get user by email
func (u *RepositoryMemory) GetUserByEmail(email string) (models.User, error) {
	var t models.User
	err := u.getByFuncOfType(func(m interface{}) bool {
		p := m.(models.User)
		return p.Email == email
	}, &t)
	return t, err
}

// Get user by id
func (u *RepositoryMemory) GetUserById(id uint) (models.User, error) {
	var m models.User
	err := u.getByIdOfType(id, &m)
	return m, err
}

// Get team by id
func (u *RepositoryMemory) GetTeam(id uint) (models.Team, error) {
	var m models.Team
	err := u.getByIdOfType(id, &m)
	return m, err
}

// Get player by id
func (u *RepositoryMemory) GetPlayer(playerId uint) (models.Player, error) {
	var m models.Player
	err := u.getByIdOfType(playerId, &m)
	return m, err
}

// Get the players of a team
func (u *RepositoryMemory) GetPlayers(teamId uint) []models.Player {
	ps := make([]models.Player, 0)
	u.getAllByFuncOfType(func(m interface{}) bool {
		p := m.(models.Player)
		return p.TeamID == teamId
	}, ps)
	return ps
}

// Get the team of a user
func (u *RepositoryMemory) GetUserTeam(user models.User) (models.Team, error) {
	var t models.Team
	err := u.getByFuncOfType(func(m interface{}) bool {
		p := m.(models.Team)
		return p.UserID == user.ID
	}, &t)
	return t, err
}

// Add a new model
func (u *RepositoryMemory) Create(model interface{}) error {
	var m interface{}
	if reflect.ValueOf(model).Kind() == reflect.Ptr {
		m = reflect.ValueOf(model).Elem().Interface()
	} else {
		m = model
	}
	u.Models = append(u.Models, m)
	return nil
}

// Update a model
func (u *RepositoryMemory) Update(model interface{}) error {
	panic("update not implemented")
}

// Delete a model
func (u *RepositoryMemory) Delete(model interface{}) error {
	panic("delete not implemented")
}

// Get all transfers
func (u *RepositoryMemory) GetTransfers() []models.Transfer {
	a := make([]models.Transfer, 0)
	u.getAllByFuncOfType(func(m interface{}) bool { return true }, a)
	return a
}

// Get transfer by id
func (u *RepositoryMemory) GetTransfer(id uint) (models.Transfer, error) {
	var m models.Transfer
	err := u.getByIdOfType(id, &m)
	return m, err
}

// Run code in a transaction (dummy)
func (u *RepositoryMemory) RunInTransaction(code func() error) error {
	return code()
}

// Delete a team
func (u *RepositoryMemory) DeleteTeam(team *models.Team) error {
	return doDeleteTeam(u, team)
}

// Delete a player
func (u *RepositoryMemory) DeletePlayer(player *models.Player) error {
	return doDeletePlayer(u, player)
}

// Get a transfer with a player
func (u *RepositoryMemory) GetTransferWithPlayer(player *models.Player) (models.Transfer, error) {
	var t models.Transfer
	err := u.getByFuncOfType(func(m interface{}) bool {
		p := m.(models.Transfer)
		return p.PlayerID == player.ID
	}, &t)
	return t, err
}

// Get model with an id and a specific type
func (u *RepositoryMemory) getByIdOfType(id uint, t interface{}) error {
	return u.getByFuncOfType(func(m interface{}) bool {
		mo := m.(gorm.Model)
		return mo.ID == id
	}, t)
}

// Get the first model of type that matches a func
func (u *RepositoryMemory) getByFuncOfType(f func(m interface{}) bool, t interface{}) error {
	r := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(reflect.ValueOf(t).Elem().Interface())), 0, 0)
	rv := reflect.New(r.Type())
	rv.Elem().Set(r)
	slicePtr := reflect.ValueOf(rv.Interface())
	sliceValuePtr := slicePtr.Elem()

	i := u.getAllByFuncOfTypeValue(f, sliceValuePtr)
	if i == 0 {
		return fmt.Errorf("not found")
	}
	reflect.ValueOf(t).Elem().Set(sliceValuePtr.Index(0))
	return nil
}

// Get all models and add them to an array
func (u *RepositoryMemory) getAllByFuncOfTypeValue(f func(m interface{}) bool, slice reflect.Value) int {
	elementType := slice.Type().Elem()
	i := 0
	for _, m := range u.Models {
		rv := reflect.ValueOf(m)
		if rv.Type().AssignableTo(elementType) && f(m) {
			slice.Set(reflect.Append(slice, rv))
			i += 1
		}
	}
	return i
}

// Get all models that match a function and are of a specific type
func (u *RepositoryMemory) getAllByFuncOfType(f func(m interface{}) bool, t interface{}) int {
	slice := reflect.ValueOf(t).Elem()
	elementType := slice.Type().Elem()
	i := 0
	for _, m := range u.Models {
		rv := reflect.ValueOf(m)
		if rv.Type().AssignableTo(elementType) && f(m) {
			slice.Set(reflect.Append(slice, rv))
			i += 1
		}
	}
	return i
}
