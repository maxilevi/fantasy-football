package models

import "gorm.io/gorm"

// Transfer DB model
type Transfer struct {
	gorm.Model
	PlayerID uint
	Player   Player
	Ask      int
}

type ShowTransfer struct {
	ID     uint       `json:"id"`
	Player ShowPlayer `json:"player"`
	Ask    int        `json:"ask"`
} //@name ShowTransfer

type UpdateTransfer struct {
	Ask int `json:"ask"`
} //@name UpdateTransfer

type CreateTransfer struct {
	PlayerID uint `json:"player_id"`
	Ask      int  `json:"ask"`
} //@name CreateTransfer
