package handlers

import (
	"../middleware"
	"../models"
	"../repos"
	"encoding/json"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/mail"
)

func AddUserRoutes(r *mux.Router, repo repos.Repository) {
	r.HandleFunc("/user", wrap(handlePostUser, repo)).Methods( "POST")
	rAuth := r.PathPrefix("/user").Subrouter()
	rAuth.Use(middleware.Auth(repo))
	rAuth.HandleFunc("", wrap(handleGetMe, repo)).Methods( "GET")
	//rAuth.HandleFunc("/{id}", wrap(handleGetUser, repo)).Methods( "GET")
}

func handleGetMe(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	user, err := getUserFromRequest(w, req)
	if err != nil {
		return
	}

	payload, err := getUserJson(user, repo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, payload)
}

func handleGetUser(w http.ResponseWriter, req *http.Request, repo repos.Repository) {
	/*user, err := getUserFromRequest(w, req, repo)

	payload, err := getUserJson(user, repo)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Internal server error")
		return
	}*/
}

func getUserJson(user models.User, repo repos.Repository) ([]byte, error) {
	team, err := repo.GetUserTeam(user)
	if err != nil {
		return nil, err
	}

	type userJson struct {
		Email string `json:"email"`
		Team uint `json:"team"`
	}

	data := userJson{
		Email: user.Email,
		Team: team.ID,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return payload, nil
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
