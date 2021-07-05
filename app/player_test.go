package app

import (
	"encoding/json"
	"gorm.io/gorm/utils/tests"
	"net/http"
	"strconv"
	"testing"
)

func TestGetPlayersFromTeam(t *testing.T) {
	setupTest()
	token, players := getTokenAndPlayerIds(t)
	for _, player := range players {
		getPlayer(t, token, player)
	}
}

func TestDeletePlayer(t *testing.T) {
	setupTest()
	token, players := getTokenAndPlayerIds(t)

	for _, player := range players {
		deletePlayer(t, token, player)
	}

	_, players = getTokenAndPlayerIds(t)
	tests.AssertEqual(t, players, make([]int, 0))
}

func TestPatchPlayer(t *testing.T) {
	setupTest()
	token, players := getTokenAndPlayerIds(t)

	player := players[0]
	payload := getPlayerPayload()

	patchPlayer(t, token, player, payload)
	getResp := getPlayer(t, token, player)

	for key, value := range getResp {
		if key == "error" {
			continue
		}
		tests.AssertEqual(t, value, payload[key])
	}
}

func TestPostPlayer(t *testing.T) {
	setupTest()
	token := getUserToken(t, "test@gmail.com")

	payload := getPlayerPayload()
	postResp := postPlayer(t, token, payload)
	getResp := getPlayer(t, token, postResp["id"].(int))
	for key, value := range getResp {
		if key == "error" {
			continue
		}
		tests.AssertEqual(t, value, payload[key])
	}
}

func getTokenAndPlayerIds(t *testing.T) (string, []int) {
	token := getUserToken(t, "test@gmail.com")
	resp, err := doGetRequest("user", token, http.StatusOK)

	if err != nil {
		t.Fatal(err)
	}

	teamId := strconv.Itoa(int(resp["team"].(float64)))
	resp, err = doGetRequest("team/"+teamId, token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	var ids []int
	for _, p := range resp["players"].([]interface {}) {
		ids = append(ids, int(p.(float64)))
	}

	return token, ids
}

func getPlayer(t *testing.T, token string, player int) map[string]interface{} {
	resp, err := doGetRequest("player/"+strconv.Itoa(player), token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func deletePlayer(t *testing.T, token string, player int) {
	playerId := strconv.Itoa(player)
	resp, err := doDeleteRequest("player/"+playerId, token, http.StatusOK)
	if err != nil || resp["error"].(bool) {
		t.Fatal(err)
	}
}

func patchPlayer(t *testing.T, token string, player int, payload map[string]interface{}) {
	playerId := strconv.Itoa(player)
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := doPatchRequest("player/"+playerId, token, body, http.StatusOK)
	if err != nil || resp["error"].(bool) {
		t.Fatal(err)
	}
}

func postPlayer(t *testing.T, token string, payload map[string]interface{}) map[string]interface{} {
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := doPostRequest("player", token, body, http.StatusOK)
	if err != nil || resp["error"].(bool) {
		t.Fatal(err)
	}
	return resp
}

func getPlayerPayload() map[string]interface{} {
	return map[string]interface{}{
		"first_name":   "test",
		"last_name":    "surname",
		"age":          123,
		"country":      "united states",
		"market_value": 10203012,
	}
}
