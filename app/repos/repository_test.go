package repos

import (
	"testing"
)

func TestRepositoryMemoryGetTeam(t *testing.T) {
	email := "test@gmail.com"
	repo := CreateRepositoryMemory()
	user, _ := repo.CreateUser(email, []byte{}, 0)
	_, err := repo.GetUserTeam(user)
	if err != nil {
		t.Error("team was not created")
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
