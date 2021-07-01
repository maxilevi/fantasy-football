package models

import (
	"github.com/jinzhu/gorm"
)

type Team struct {
	gorm.Model
	Name string
	Country string
	Players []Player
}

func (team *Team) MarketValue() int32 {
	var sum int32
	for _, player := range team.Players {
		sum += player.MarketValue
	}
	return sum
}
