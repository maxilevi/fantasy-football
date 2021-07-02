package repos

import (
	"../models"
	"testing"
)

func TestRepositoryMemory_CreateUser(t *testing.T) {
	email := "test@gmail.com"
	repo := &RepositoryMemory{Users: make([]models.User, 0)}
	repo.CreateUser(email, []byte{}, 0)
	if len(repo.Users) != 1 {
		t.Error("user was not created")
	}
}

func TestRepositoryMemory_GetUser(t *testing.T) {
	email := "test@gmail.com"
	repo := &RepositoryMemory{Users: make([]models.User, 0)}
	repo.CreateUser(email, []byte{}, 0)

	var user models.User
	err := repo.GetUser(email, &user)
	if err != nil {
		t.Error(err)
	}
}