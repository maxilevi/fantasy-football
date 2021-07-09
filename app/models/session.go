package models


type CreateSession struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}