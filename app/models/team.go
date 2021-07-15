package models

import (
	"github.com/Pallinder/go-randomdata"
	"github.com/jinzhu/gorm"
)

const (
	teamSize          = 20
	goalKeeperCount   = 3
	defenderCount     = 6
	midfielderCount   = 6
	attackerCount     = 5
	DefaultTeamBudget = 5000000
)

// Team DB model
type Team struct {
	gorm.Model
	Name    string
	Country string
	Budget  int
	UserID  uint
	User    User
}

// Create a team with a random name, default budget and random players
func RandomTeam() (Team, []Player) {
	team := Team{
		Name:    randomdata.SillyName(),
		Country: randomdata.Country(randomdata.FullCountry),
		Budget:  DefaultTeamBudget,
	}
	players := make([]Player, teamSize)
	i := 0
	for j := i; i < j+goalKeeperCount; i++ {
		players[i] = RandomPlayer(goalkeeper)
	}
	for j := i; i < j+defenderCount; i++ {
		players[i] = RandomPlayer(defender)
	}
	for j := i; i < j+midfielderCount; i++ {
		players[i] = RandomPlayer(midfielder)
	}
	for j := i; i < j+attackerCount; i++ {
		players[i] = RandomPlayer(attacker)
	}
	return team, players
}

type ShowTeam struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Country     string       `json:"country"`
	Budget      int          `json:"budget"`
	MarketValue int          `json:"market_value"`
	Players     []ShowPlayer `json:"players"`
} //@name ShowTeam

type CreateTeam struct {
	Owner   int    `json:"owner" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Country string `json:"country" binding:"required"`
	Budget  int    `json:"budget" binding:"required"`
} //@name CreateTeam

type UpdateTeam struct {
	Name    string `json:"name"`
	Country string `json:"country"`
	Budget  int    `json:"budget"`
} //@name UpdateTeam
