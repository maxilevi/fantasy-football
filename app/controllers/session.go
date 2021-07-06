package controllers

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

type SessionController struct {
	Repo repos.Repository
}

func (c *SessionController) AddRoutes(r *mux.Router) {
	r.HandleFunc("/session", c.handlePostSession).Methods("POST")
}

// Handles creating a new session
func (c *SessionController) handlePostSession(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	type sessionCreation struct {
		Email    string
		Password string
	}

	var t sessionCreation
	err := decoder.Decode(&t)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Incorrect body parameters")
		return
	}
	user, err := c.loginAndGetUser(t.Email, t.Password)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	token, err := c.createToken(user)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeResponse(w, http.StatusOK, []byte(`{"error": false, "token": "`+token+`"}`))
}

// Tries to login the user and returns it if successful
func (c *SessionController) loginAndGetUser(email, password string) (*models.User, error) {
	user, err := c.Repo.GetUserByEmail(email)
	if err != nil {
		return &user, err
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password))
	if err != nil {
		return &user, err
	}

	return &user, nil
}

// Creates a JWT token for a specific user
func (c *SessionController) createToken(user *models.User) (string, error) {
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
