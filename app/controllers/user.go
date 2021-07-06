package controllers

import (
	"../middleware"
	"../models"
	"../repos"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"net/mail"
)

type UserController struct {
	Repo repos.Repository
}

func (c *UserController) AddRoutes(r *mux.Router) {

	r.HandleFunc("/user", c.handlePostUser).Methods("POST")

	rAdmin := r.PathPrefix("/user").Subrouter()
	rAdmin.Use(middleware.Auth(c.Repo))
	rAdmin.Use(middleware.Admin)
	rAdmin.HandleFunc("/{id}", c.handleGetUser).Methods("GET")
	rAdmin.HandleFunc("/{id}", c.handleDeleteUser).Methods("DELETE")
	rAdmin.HandleFunc("/{id}", c.handlePatchUser).Methods("PATCH")

	rAuth := r.PathPrefix("/user").Subrouter()
	rAuth.Use(middleware.Auth(c.Repo))
	rAuth.HandleFunc("", c.handleGetMe).Methods("GET")
}

// Handles GET request to the user resource when no ID is provided
func (c *UserController) handleGetMe(w http.ResponseWriter, req *http.Request) {
	user, err := getAuthenticatedUserFromRequest(w, req)
	if err != nil {
		return
	}

	payload, err := c.makeUserJson(user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, payload)
}

// Handles GET request to the user resource
func (c *UserController) handleGetUser(w http.ResponseWriter, req *http.Request) {
	user, err := c.getUserFromRequest(w, req)
	if err != nil {
		return
	}

	resp, err := c.makeUserJson(user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeResponse(w, http.StatusOK, resp)
}

// Handles DELETE requests to the user's resource
func (c *UserController) handleDeleteUser(w http.ResponseWriter, req *http.Request) {
	user, err := c.getUserFromRequest(w, req)
	if err != nil {
		return
	}

	err = c.Repo.Delete(user)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}
	writeResponse(w, http.StatusOK, noError())
}

// Handles PATCH requests to the user's resource
func (c *UserController) handlePatchUser(w http.ResponseWriter, req *http.Request) {
	/*user, err := c.getAuthenticatedUserFromRequest(w, req)
	if err != nil {
		return
	}

	payload, err := getUserJson(user, Repo)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}*/
}

// Handles POST request to the user resource
func (c *UserController) handlePostUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	type userRegistration struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var t userRegistration
	err := decoder.Decode(&t)
	if err != nil {
		writeErrorResponse(w, http.StatusBadRequest, "Incorrect body parameters")
		return
	}
	if !c.validEmail(t.Email) {
		writeErrorResponse(w, http.StatusBadRequest, "Invalid email")
		return
	}
	if c.emailExists(t.Email) {
		writeErrorResponse(w, http.StatusBadRequest, "Provided email is already registered")
		return
	}
	if !c.validPassword(t.Password) {
		writeErrorResponse(w, http.StatusBadRequest, "Password needs a minimum of at least 8 characters")
		return
	}

	log.Println("Registering a new user...")

	user, err := c.registerUser(t.Email, t.Password)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	writeResponse(w, http.StatusOK, []byte(fmt.Sprintf(`{"error": false, "id": %v}`, user.ID)))
}

// Return a json representation from a given user
func (c *UserController) makeUserJson(user models.User) ([]byte, error) {
	team, err := c.Repo.GetUserTeam(user)
	if err != nil {
		return nil, err
	}

	type userJson struct {
		Email string `json:"email"`
		Team  uint   `json:"team"`
	}

	data := userJson{
		Email: user.Email,
		Team:  team.ID,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// Returns a bool to check if the password is valid
func (c *UserController) validPassword(password string) bool {
	return len(password) >= 8
}

// Returns a bool to check if the email is valid
func (c *UserController) validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// Return a bool to see if the email is already registered
func (c *UserController) emailExists(email string) bool {
	_, err := c.Repo.GetUserByEmail(email)
	return err == nil
}

// Registers a new user with the given credentials
func (c *UserController) registerUser(email, password string) (models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	return c.Repo.CreateUser(email, hashedPassword, 0)
}

// Parse a user from the request parameters or return an error if not found
func (c *UserController) getUserFromRequest(w http.ResponseWriter, req *http.Request) (models.User, error) {
	id, err := parseIdFromRequest(w, req)
	if err != nil {
		return models.User{}, err
	}

	user, err := c.Repo.GetUserById(id)
	if err != nil {
		writeErrorResponse(w, http.StatusNotFound, "Not found")
		return models.User{}, err
	}
	return user, err
}
