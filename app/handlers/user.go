package handlers

import (
	"../models"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/mail"
	"../repos"
)

func AddUserRoutes(r *mux.Router, repo repos.Repository) {
	r.HandleFunc("/users", wrap(handlePostUser, repo)).Methods( "POST")
}

type userRegistration struct {
	Email string
	Password string
}

func handlePostUser(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	decoder := json.NewDecoder(req.Body)
	var t userRegistration
	err := decoder.Decode(&t)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Incorrect body parameters")
		return
	}
	if !validEmail(t.Email) {
		writeError(w, http.StatusBadRequest, "Invalid email")
		return
	}
	if emailExists(t.Email, repo) {
		writeError(w, http.StatusBadRequest, "Provided email is already registered")
		return
	}
	if !validPassword(t.Password) {
		writeError(w, http.StatusBadRequest, "Password needs a minimum of at least 8 characters")
		return
	}

	log.Println("Registering a new user...")

	err = registerUser(t, repo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, 200, []byte(`{"error": false}`))
}

func validPassword(password string) bool {
	return len(password) >= 8
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func emailExists(email string, repo repos.Repository) bool {
	var user models.User
	err := repo.GetUser(email, &user)
	return err == nil
}

func registerUser(reg userRegistration, repo repos.Repository) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	repo.CreateUser(reg.Email, hashedPassword, 0)
	return nil
}