package models

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/jinzhu/gorm"
)

const teamSize = 20

type Team struct {
	gorm.Model
	Name string
	Country string
	Players []Player
	Owner User
}

func (team *Team) MarketValue() int32 {
	var sum int32
	for _, player := range team.Players {
		sum += player.MarketValue
	}
	return sum
}

func RandomTeam() Team {
	players := make([]Player, teamSize)
	for i := 0; i < len(players); i++ {
		players[i] = RandomPlayer()
	}
	team := Team{
		Players: players,
		Name: randomdata.SillyName(),
		Country: randomdata.Country(randomdata.FullCountry),
	}
	return team
}