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
	OwnerID uint
	Owner User
}

func RandomTeam() (Team, []Player) {
	team := Team{
		Name: randomdata.SillyName(),
		Country: randomdata.Country(randomdata.FullCountry),
	}
	players := make([]Player, teamSize)
	for i := 0; i < len(players); i++ {
		players[i] = RandomPlayer()
	}
	return team, players
}