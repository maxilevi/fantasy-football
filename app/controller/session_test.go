package controller

import (
	"../middleware"
	"../models"
	"../repos"
	"fmt"
	"github.com/golang-jwt/jwt"
	"testing"
)

func TestGetUser(t *testing.T) {
	db := repos.CreateRepositoryMemory()
	c := Controller{Repo: db}
	email := "test@gmail.com"
	pass := "hello123"
	_, err := c.registerUser(email, pass)
	if err != nil {
		t.Error(err)
		return
	}

	user2, err := c.loginAndGetUser(email, pass)
	if err != nil {
		t.Error(err)
		return
	}
	if user2.Email != email {
		t.Error("wrong user")
	}

}

func TestCreateToken(t *testing.T) {
	c := Controller{Repo: repos.CreateRepositoryMemory()}
	email := "test@gmail.com"
	tokenString, err := c.createToken(&models.User{Email: email})
	if err != nil {
		t.Error("failed to create token")
	}

	claims := &middleware.AuthClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return middleware.JWTKey, nil
	})

	if err != nil || !token.Valid {
		t.Error("invalid token")
	}

	if claims.Email != email {
		t.Error("token has invalid email")
	}
}
