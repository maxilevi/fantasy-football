package models

import "gorm.io/gorm"

type Transfer struct {
	gorm.Model
	PlayerID int
	Player Player
	Ask int
}