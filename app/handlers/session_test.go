package handlers

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"testing"
	"../repos"
	"../models"
	"../middleware"
)

func TestGetUser(t *testing.T) {
	email := "test@gmail.com"
	pass := "hello123"
	repo := &repos.RepositoryMemory{Users: make([]models.User, 0)}
	err := registerUser(userRegistration{Email: email, Password: pass}, repo)
	if err != nil {
		t.Error()
		return
	}
	user, err := getUser(sessionCreation{Email: email, Password: pass}, repo)
	if err != nil {
		t.Error()
		return
	}
	if user.Email != email {
		t.Error("wrong user")
	}

}

func TestCreateToken(t *testing.T) {
	email := "test@gmail.com"
	tokenString, err := createToken(&models.User{Email: email})
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