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

// DB player model
type Player struct {
	gorm.Model
	FirstName   string
	LastName    string
	Country     string
	Age         int
	MarketValue int32
	Position    int
	TeamID      uint
	Team        Team
}

// Create a player with random characteristics
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

// Get a random age between 18 and 40
func randomAge() int {
	return rand.Intn(40-18) + 18
}

type BasePlayer struct {
	FirstName   string `json:"first_name" example:"Audrey"`
	LastName    string `json:"last_name" example:"Hepburn"`
	Country     string `json:"country" example:"Germany"`
	Age         int    `json:"age" example:"25"`
	MarketValue int32  `json:"market_value" example:"25000"`
	// This is the position identifier 0 for goalkeeper, 1 for defender, 2 for goalkeeper, 3 for attacker
	Position int `json:"position" example:"1" validate:"min=0,max=3" minimum:"0" maximum:"3"`
}

type ShowPlayer struct {
	BasePlayer
	ID uint `json:"id"`
} //@name ShowPlayer

type CreatePlayer struct {
	FirstName   string `json:"first_name" example:"Audrey" binding:"required"`
	LastName    string `json:"last_name" example:"Hepburn" binding:"required"`
	Country     string `json:"country" example:"Germany" binding:"required"`
	Age         int    `json:"age" example:"25" binding:"required"`
	MarketValue int32  `json:"market_value" example:"25000" binding:"required"`
	// This is the position identifier 0 for goalkeeper, 1 for defender, 2 for goalkeeper, 3 for attacker
	Position int `json:"position" example:"1" binding:"required" validate:"min=0,max=3" minimum:"0" maximum:"3"`
} //@name CreatePlayer

type UpdatePlayer struct {
	BasePlayer
	Team int `json:"team"`
} //@name UpdatePlayer
