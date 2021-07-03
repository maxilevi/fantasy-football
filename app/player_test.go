package app

import (
	"strconv"
	"testing"
)

func TestGetPlayersFromTeam(t *testing.T) {
	token, players := getTokenAndPlayerIds(t)
	for _, player := range players {
		_, err := doGetRequest("player/" + strconv.Itoa(player), token)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Cleanup(func() { truncateDb() })
}

func TestDeletePlayer(t *testing.T) {
	token, players := getTokenAndPlayerIds(t)

	t.Cleanup(func() { truncateDb() })
}

func TestPatchPlayer(t *testing.T) {
	token, players := getTokenAndPlayerIds(t)

	t.Cleanup(func() { truncateDb() })
}

func TestPostPlayer(t *testing.T) {
	token := getUserToken(t, "test@gmail.com")

	t.Cleanup(func() { truncateDb() })
}

func getTokenAndPlayerIds(t *testing.T) (string, []int) {
	token := getUserToken(t, "test@gmail.com")
	resp, err := doGetRequest("user", token)

	if err != nil {
		t.Fatal(err)
	}

	teamId := strconv.Itoa(int(resp["team"].(float64)))
	resp, err = doGetRequest("team/"+teamId, token)
	if err != nil {
		t.Fatal(err)
	}
	return token, resp["players"].([]int)
}