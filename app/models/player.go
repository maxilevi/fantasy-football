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


type ShowPlayer struct {
	ID 			uint    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Country     string `json:"country"`
	Age         int    `json:"age"`
	MarketValue int    `json:"market_value"`
	Position    int    `json:"position"`
}

type CreatePlayer struct {
	FirstName   string
	LastName    string
	Country     string
	Age         int
	MarketValue int32
	Position	int
	TeamID      uint
}

type UpdatePlayer struct {
	Team 		int    `json:"team"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Country     string `json:"country"`
	Age         int    `json:"age"`
	MarketValue int    `json:"market_value"`
	Position    int    `json:"position"`
}