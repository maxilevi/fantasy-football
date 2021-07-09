package models

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/jinzhu/gorm"
	"math/rand"
)

const (
	goalkeeper = iota
	defender
	midfielder
	attacker
)

type Player struct {
	gorm.Model
	FirstName   string
	LastName    string
	Country     string
	Age         int
	MarketValue int32
	Position	int
	TeamID      uint
	Team        Team
}

func RandomPlayer(position int) Player {
	return Player{
		MarketValue: 1000000,
		FirstName:   randomdata.FirstName(randomdata.RandomGender),
		LastName:    randomdata.LastName(),
		Age:         randomAge(),
		Country:     randomdata.Country(randomdata.FullCountry),
		Position:    position,
	}
}

func randomAge() int {
	return rand.Intn(40-18) + 18
}

type BasePlayer struct {
	FirstName   string `json:"first_name" example:"Audrey"`
	LastName    string `json:"last_name" example:"Hepburn"`
	Country     string `json:"country" example:"Germany"`
	Age         int `json:"age" example:"25"`
	MarketValue int32 `json:"market_value" example:"25000"`
	Position	int `json:"position" example:"1"`
}

type ShowPlayer struct {
	BasePlayer
	ID 			uint    `json:"id"`
}

type CreatePlayer struct {
	BasePlayer
	Team      uint `json:"team"`
}

type UpdatePlayer struct {
	BasePlayer
	Team 		int    `json:"team"`
}