package models

import (
	"github.com/jinzhu/gorm"
	"math/rand"
	"github.com/Pallinder/go-randomdata"
)

type Player struct {
	gorm.Model
	FirstName string
	LastName string
	Country string
	Age int
	MarketValue int32
	TeamID uint
	Team Team
}

func RandomPlayer() Player {
	return Player{
		MarketValue: 1000000,
		FirstName: randomdata.FirstName(randomdata.RandomGender),
		LastName: randomdata.LastName(),
		Age: randomAge(),
		Country: randomdata.Country(randomdata.FullCountry),
	}
}


func randomAge() int {
	return rand.Intn(40 - 18) + 18
}
