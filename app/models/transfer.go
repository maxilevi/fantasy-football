package models

import "gorm.io/gorm"

type Transfer struct {
	gorm.Model
	PlayerID int
	Player   Player
	Ask      int
}

type ShowTransfer struct {
	ID     uint       `json:"id"`
	Player ShowPlayer `json:"player"`
	Ask    int        `json:"ask"`
}

type UpdateTransfer struct {
	PlayerID uint `json:"player_id"`
	Ask      int  `json:"ask"`
}

type CreateTransfer struct {
	PlayerID uint `json:"player_id"`
	Ask      int  `json:"ask"`
}
