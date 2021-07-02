package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"../models"
	"../repos"
)

func AddSessionRoutes(r *mux.Router, repo *repos.Repository) {
	r.HandleFunc("/session", wrap(handlePostSession, repo)).Methods( "POST")
}

type sessionCreation struct {
	Email string
	Password string
}

func handlePostSession(w http.ResponseWriter, req *http.Request, repo *repos.Repository) {
	decoder := json.NewDecoder(req.Body)
	var t sessionCreation
	err := decoder.Decode(&t)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Incorrect body parameters")
		return
	}
	user, err := getUser(t, repo)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid email or password")
		return
	}
	session, err := createSession(user, repo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

}

func getUser(params sessionCreation, repo *repos.Repository) (models.User, error) {
	var user models.User
	err := repo.GetUser(&models.User{Email: params.Email})
	if err != nil {
		return user, err
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(params.Password))
	if err != nil {
		return user, err
	}
	return user, nil
}

func createSession(user *models.User, repo *repos.Repository) error {
	//TODO
}