package handlers

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)
import "../repos"
import "../models"

func TestValidPassword(t *testing.T) {
	if validPassword("1234") || validPassword("1234567") || validPassword("")  {
		t.Error("Invalid password was marked as valid")
	}
	if !(validEmail("12345678") && validEmail("12345678890") ) {
		t.Error("Valid password was marked as invalid")
	}
}

func TestEmail(t *testing.T) {
	if validEmail("test") || validEmail("test@") || validEmail("@gmail") || validEmail("") {
		t.Error("Invalid email was marked as valid")
	}
	if !(validEmail("test@gmail.com") && validEmail("a.a@a.a") && validEmail("test@something.com")) {
		t.Error("Valid email was marked as invalid")
	}
}

func TestHasEmail(t *testing.T) {
	email := "test@gmail.com"
	repo := repos.RepositoryMemory{Users: make([]models.User, 10)}
	repo.CreateUser(email, []byte{}, 0)

	if !emailExists(email, repo) {
		t.Error()
	}
	if emailExists("nosuchemail", repo) {
		t.Error()
	}
}

func TestRegisterUser(t *testing.T) {
	email := "test@gmail.com"
	pass := "hello123"
	repo := repos.RepositoryMemory{Users: make([]models.User, 10)}
	err := registerUser(userRegistration{Email: email, Password: pass}, repo)
	if err != nil {
		t.Error()
	}
	var user models.User
	err = repo.GetUser("test@gmail.com", &user)
	if err != nil || user.Email != email || user.PermissionLevel != 0 {
		t.Error()
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pass))
	if err != nil {
		t.Error("password hash does not match")
	}
}