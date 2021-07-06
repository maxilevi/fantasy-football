package handlers

import (
	"../middleware"
	"../models"
	"../repos"
	"encoding/json"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

func AddSessionRoutes(r *mux.Router, repo repos.Repository) {
	r.HandleFunc("/session", wrap(handlePostSession, repo)).Methods("POST")
}

type sessionCreation struct {
	Email    string
	Password string
}

func handlePostSession(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	decoder := json.NewDecoder(req.Body)
	var t sessionCreation
	err := decoder.Decode(&t)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Incorrect body parameters")
		return
	}
	user, err := getUser(t, repo)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	token, err := createToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeResponse(w, 200, []byte(`{"error": false, "token": "`+token+`"}`))
}

func getUser(params sessionCreation, repo repos.Repository) (*models.User, error) {
	var user models.User
	err := repo.GetUser(params.Email, &user)
	if err != nil {
		return &user, err
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(params.Password))
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func createToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middleware.AuthClaims{
		Email: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(middleware.JWTKey)
}
