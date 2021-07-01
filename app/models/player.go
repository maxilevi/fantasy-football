package models

import (
"github.com/jinzhu/gorm"
)

type Player struct {
	gorm.Model
	FirstName string
	LastName string
	Country string
	Age string
	MarketValue int32
}

func RandomPlayer() Player {
	return Player{

	}
}