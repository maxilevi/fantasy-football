package repos

import (
	"../models"
	"testing"
)

func TestRepositoryMemoryCreateUser(t *testing.T) {
	email := "test@gmail.com"
	repo := CreateRepositoryMemory()
	repo.CreateUser(email, []byte{}, 0)
	if len(repo.Users) != 1 {
		t.Error("user was not created")
	}
}

func TestRepositoryMemoryGetUser(t *testing.T) {
	email := "test@gmail.com"
	repo := CreateRepositoryMemory()
	repo.CreateUser(email, []byte{}, 0)

	var user models.User
	err := repo.GetUser(email, &user)
	if err != nil {
		t.Error(err)
	}
}
