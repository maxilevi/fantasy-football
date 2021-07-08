package repos

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)
import "../models"

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
	RunInTransaction(code func () error) error
}

func doCreateUser(u Repository, email string, hash []byte, permission int) (models.User, error) {
	user := models.User{
		Email:           email,
		PasswordHash:    hash,
		PermissionLevel: permission,
	}
	return user, u.RunInTransaction(func () error {
		err := u.Create(&user)
		if err != nil {
			return err
		}

		team, players := models.RandomTeam()
		team.OwnerID = user.ID
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

type RepositorySQL struct {
	Db *gorm.DB
}

func (u RepositorySQL) CreateUser(email string, hash []byte, permission int) (models.User, error) {
	return doCreateUser(u, email, hash, permission)
}

func (u RepositorySQL) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	res := u.Db.Where(models.User{Email: email}).First(&user)
	return user, res.Error
}

func (u RepositorySQL) GetUserById(id uint) (models.User, error) {
	var user models.User
	res := u.Db.First(&user, id)
	return user, res.Error
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

func (u RepositorySQL) Create(model interface{}) error {
	res := u.Db.Save(model)
	return res.Error
}

func (u RepositorySQL) Update(model interface{}) error {
	res := u.Db.Save(model)
	return res.Error
}

func (u RepositorySQL) Delete(model interface{}) error {
	res := u.Db.Delete(model)
	return res.Error
}

func (u RepositorySQL) GetUserTeam(user models.User) (models.Team, error) {
	var team models.Team
	res := u.Db.Where(&models.Team{OwnerID: user.ID}).Find(&team)
	return team, res.Error
}

func (u RepositorySQL) GetTransfers() []models.Transfer {
	var transfers []models.Transfer
	u.Db.Where("1 = 1").Find(&transfers)
	return transfers
}

func (u RepositorySQL) GetTransfer(id uint) (models.Transfer, error) {
	var transfer models.Transfer
	res := u.Db.Find(&transfer, id)
	return transfer, res.Error
}

func (u RepositorySQL) RunInTransaction(code func () error) error {
	tx := u.Db.Begin()
	err := code()
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

type RepositoryMemory struct {
	Models   []interface{}
}

func CreateRepositoryMemory() *RepositoryMemory {
	return &RepositoryMemory{
		Models:   make([]interface{}, 0),
	}
}

func (u *RepositoryMemory) CreateUser(email string, hash []byte, permission int) (models.User, error) {
	return doCreateUser(u, email, hash, permission)
}

func (u *RepositoryMemory) GetUserByEmail(email string) (models.User, error) {
	var t models.User
	err := u.getByFuncOfType(func(m interface{}) bool {
		p := m.(models.User)
		return p.Email == email
	}, &t)
	return t, err
}

func (u *RepositoryMemory) GetUserById(id uint) (models.User, error) {
	var m models.User
	err := u.getByIdOfType(id, &m)
	return m, err
}

func (u *RepositoryMemory) GetTeam(id uint) (models.Team, error) {
	var m models.Team
	err := u.getByIdOfType(id, &m)
	return m, err
}

func (u *RepositoryMemory) GetPlayer(playerId uint) (models.Player, error) {
	var m models.Player
	err := u.getByIdOfType(playerId, &m)
	return m, err
}

func (u *RepositoryMemory) GetPlayers(teamId uint) []models.Player {
	ps := make([]models.Player, 0)
	u.getAllByFuncOfType(func(m interface{}) bool {
		p := m.(models.Player)
		return p.TeamID == teamId
	}, ps)
	return ps
}

func (u *RepositoryMemory) GetUserTeam(user models.User) (models.Team, error) {
	var t models.Team
	err := u.getByFuncOfType(func(m interface{}) bool {
		p := m.(models.Team)
		return p.OwnerID == user.ID
	}, &t)
	return t, err
}

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

func (u *RepositoryMemory) Update(model interface{}) error {
	panic("delete not implemented")
}

func (u *RepositoryMemory) Delete(model interface{}) error {
	panic("delete not implemented")
}

func (u *RepositoryMemory) GetTransfers() []models.Transfer {
	a := make([]models.Transfer, 0)
	u.getAllByFuncOfType(func(m interface{}) bool {return true}, a)
	return a
}

func (u *RepositoryMemory) GetTransfer(id uint) (models.Transfer, error) {
	var m models.Transfer
	err := u.getByIdOfType(id, &m)
	return m, err
}

func (u *RepositoryMemory) RunInTransaction(code func() error) error {
	return code()
}

func (u *RepositoryMemory) getByIdOfType(id uint, t interface{}) error {
	return u.getByFuncOfType(func (m interface{}) bool {
		mo := m.(gorm.Model)
		return mo.ID == id
	}, t)
}

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
