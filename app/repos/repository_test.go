package repos

import (
	"testing"
)

func TestRepositoryMemoryCreateUser(t *testing.T) {
	email := "test@gmail.com"
	repo := CreateRepositoryMemory()
	_, _ = repo.CreateUser(email, []byte{}, 0)
	if len(repo.Models) != 1 {
		t.Error("user was not created")
	}
}

func TestRepositoryMemoryGetUser(t *testing.T) {
	email := "test@gmail.com"
	repo := CreateRepositoryMemory()
	_, _ = repo.CreateUser(email, []byte{}, 0)

	_, err := repo.GetUserByEmail(email)
	if err != nil {
		t.Error(err)
	}
}
