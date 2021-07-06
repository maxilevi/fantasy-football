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
