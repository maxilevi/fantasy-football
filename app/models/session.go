package models

type CreateSession struct {
	Email    string `json:"email"`
	Password string `json:"password"`
} //@name Credentials

type SessionToken struct {
	Token string `json:"token"`
} //@name Token
