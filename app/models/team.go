package models

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/jinzhu/gorm"
)

const (
	teamSize = 20
	goalKeeperCount = 3
	defenderCount = 6
	midfielderCount = 6
	attackerCount = 5
)

type Team struct {
	gorm.Model
	Name    string
	Country string
	Budget  int
	OwnerID uint
	Owner   User
}

func RandomTeam() (Team, []Player) {
	team := Team{
		Name:    randomdata.SillyName(),
		Country: randomdata.Country(randomdata.FullCountry),
		Budget:  5000000,
	}
	players := make([]Player, teamSize)
	i := 0
	for j := i; i < j + goalKeeperCount; i++ {
		players[i] = RandomPlayer(goalkeeper)
	}
	for j := i; i < j + defenderCount; i++ {
		players[i] = RandomPlayer(defender)
	}
	for j := i; i < j + midfielderCount; i++ {
		players[i] = RandomPlayer(midfielder)
	}
	for j := i; i < j + attackerCount; i++ {
		players[i] = RandomPlayer(attacker)
	}
	return team, players
}


type ShowTeam struct {
	ID uint `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Budget  int `json:"budget"`
	MarketValue  int `json:"market_value"`
	Players []ShowPlayer `json:"players"`
}

type CreateTeam struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Budget  int `json:"budget"`
}

type UpdateTeam struct {
	CreateTeam
}