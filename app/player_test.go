package app

import (
	"gorm.io/gorm/utils/tests"
	"net/http"
	"strconv"
	"testing"
)

func TestGetPlayersFromTeam(t *testing.T) {
	setupTest()
	token, players := getTokenAndPlayerIds(t, false)
	for _, player := range players {
		getPlayer(t, token, player)
	}
}

func TestDeletePlayer(t *testing.T) {
	setupTest()
	token, players := getTokenAndPlayerIds(t, true)

	for _, player := range players {
		deletePlayer(t, token, player)
	}

	players = getPlayersFromToken(t, token)
	tests.AssertEqual(t, players, make([]int, 0))
}

func TestPatchPlayer(t *testing.T) {
	setupTest()
	token, players := getTokenAndPlayerIds(t, true)

	player := players[0]
	payload := getPlayerPayload()

	patchPlayer(t, token, player, payload)
	getResp := getPlayer(t, token, player)

	for key, value := range getResp {
		if key == "error" || key == "team" || key == "id" {
			continue
		}
		tests.AssertEqual(t, value, payload[key])
	}
}

func TestPostPlayer(t *testing.T) {
	setupTest()
	token := getAdminUserToken(t, "test@gmail.com")


	payload := getPlayerPayload()
	payload["team"] = getTeamIdFromUser(t, token)
	postResp := postPlayer(t, token, payload)
	getResp := getPlayer(t, token, int(postResp["id"].(float64)))
	for key, value := range getResp {
		if key == "error" || key == "id" {
			continue
		}
		tests.AssertEqual(t, value, payload[key])
	}
}

func getTeamIdFromUser(t *testing.T, token string) int {
	resp, err := doGetRequest("user/me", token, http.StatusOK)

	if err != nil {
		t.Fatal(err)
	}
	team := resp["team"].(map[string]interface{})
	return int(team["id"].(float64))
}

func getTokenAndPlayerIds(t *testing.T, admin bool) (string, []int) {
	var token string
	if admin {
		token = getAdminUserToken(t, "test@gmail.com")
	} else {
		token = getUserToken(t, "test@gmail.com")
	}
	ids := getPlayersFromToken(t, token)

	return token, ids
}

func getPlayersFromToken(t *testing.T, token string) []int {
	resp, err := doGetRequest("user/me", token, http.StatusOK)

	if err != nil {
		t.Fatal(err)
	}

	team := resp["team"].(map[string]interface{})

	var ids []int
	for _, p := range team["players"].([]interface{}) {
		m := p.(map[string]interface{})
		ids = append(ids, int(m["id"].(float64)))
	}
	return ids
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
	if err != nil || int(resp["code"].(float64)) != 200 {
		t.Fatal(err)
	}
}

func patchPlayer(t *testing.T, token string, player int, payload map[string]interface{}) {
	playerId := strconv.Itoa(player)
	resp, err := doPatchRequest("player/"+playerId, token, payload, http.StatusOK)
	if err != nil || int(resp["code"].(float64)) != 200 {
		t.Fatal(err)
	}
}

func postPlayer(t *testing.T, token string, payload map[string]interface{}) map[string]interface{} {
	resp, err := doPostRequest("player", token, payload, http.StatusOK)
	if err != nil || int(resp["code"].(float64)) != 200 {
		t.Fatal(err)
	}
	return resp
}

func getPlayerPayload() map[string]interface{} {
	return map[string]interface{}{
		"first_name":   "test",
		"last_name":    "surname",
		"age":          123,
		"position": 	1,
		"country":      "united states",
		"market_value": 10203012,
	}
}
