package controllers

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)
import "../repos"

func TestValidPassword(t *testing.T) {
	c := UserController{Repo: repos.CreateRepositoryMemory()}
	if c.validPassword("1234") || c.validPassword("1234567") || c.validPassword("") {
		t.Error("Invalid password was marked as valid")
	}
	if !(c.validPassword("12345678") && c.validPassword("12345678890")) {
		t.Error("Valid password was marked as invalid")
	}
}

func TestEmail(t *testing.T) {
	c := UserController{Repo: repos.CreateRepositoryMemory()}
	if c.validEmail("test") || c.validEmail("test@") || c.validEmail("@gmail") || c.validEmail("") {
		t.Error("Invalid email was marked as valid")
	}
	if !(c.validEmail("test@gmail.com") && c.validEmail("a.a@a.a") && c.validEmail("test@something.com")) {
		t.Error("Valid email was marked as invalid")
	}
}

func TestHasEmail(t *testing.T) {
	c := UserController{Repo: repos.CreateRepositoryMemory()}
	email := "test@gmail.com"
	repo := repos.CreateRepositoryMemory()
	repo.CreateUser(email, []byte{}, 0)

	if !c.emailExists(email) {
		t.Error()
	}
	if c.emailExists("nosuchemail") {
		t.Error()
	}
}

func TestRegisterUser(t *testing.T) {
	c := UserController{Repo: repos.CreateRepositoryMemory()}
	email := "test@gmail.com"
	pass := "hello123"
	user, err := c.registerUser(email, pass)
	if err != nil {
		t.Error(err)
	}
	if user.Email != email || user.PermissionLevel != 0 {
		t.Error("user does not match")
	}
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(pass))
	if err != nil {
		t.Error("password hash does not match")
	}
}
