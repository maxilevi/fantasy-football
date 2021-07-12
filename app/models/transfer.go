package models

import "gorm.io/gorm"

type Transfer struct {
	gorm.Model
	PlayerID uint
	Player   Player
	Ask      int
	SellerID uint
	Seller Team
}

type ShowTransfer struct {
	ID     uint       `json:"id"`
	Player ShowPlayer `json:"player"`
	Ask    int        `json:"ask"`
}

type UpdateTransfer struct {
	Ask int `json:"ask"`
}

type CreateTransfer struct {
	PlayerID uint `json:"player_id"`
	Ask      int  `json:"ask"`
}
