package handlers

import (
	"../models"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/mail"
)

func AddUserRoutes(r *mux.Router, db *gorm.DB) {
	r.HandleFunc("/users", wrap(handlePostUser, db)).Methods( "POST")
}

type userRegistration struct {
	Email string
	Password string
}

func handlePostUser(w http.ResponseWriter, req *http.Request, db *gorm.DB) {
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
	if emailExists(t.Email, db) {
		writeError(w, http.StatusBadRequest, "Provided email is already registered")
		return
	}
	if !validPassword(t.Password) {
		writeError(w, http.StatusBadRequest, "Password needs a minimum of at least 8 characters")
		return
	}

	log.Println("Registering a new user...")

	err = registerUser(t)
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

func emailExists(email string, db *gorm.DB) bool {
	var user models.User
	result := db.Where(&models.User{Email: email}).First(&user)
	return result.Error != nil
}

func registerUser(reg userRegistration) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := models.User{}
	user.Email = reg.Email
	user.PasswordHash = hashedPassword
	user.PermissionLevel = 0
	return nil
}