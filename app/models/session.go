package models

type CreateSession struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
} //@name Credentials

type SessionToken struct {
	Token string `json:"token"`
} //@name Token
