package app

import (
	"gorm.io/gorm/utils/tests"
	"net/http"
	"strconv"
	"testing"
)

func TestQueryUserAndTeamInformation(t *testing.T) {
	setupTest()
	token := getUserToken(t, "test@gmail.com")
	resp, err := doGetRequest("users/me", token, http.StatusOK)

	if err != nil || resp["email"].(string) != "test@gmail.com" {
		t.Fatal(err)
	}

	team := resp["team"].(map[string]interface{})
	teamId := strconv.Itoa(int(team["id"].(float64)))
	resp, err = doGetRequest("teams/"+teamId, token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPatchTeam(t *testing.T) {
	setupTest()
	token := getUserToken(t, "test@gmail.com")
	resp, err := doGetRequest("users/me", token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	team := resp["team"].(map[string]interface{})
	res := "teams/" + strconv.Itoa(int(team["id"].(float64)))
	resp, err = doPatchRequest(res, token, map[string]interface{}{
		"name":    "New name",
		"country": "New country",
	}, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = doGetRequest(res, token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	tests.AssertEqual(t, resp["name"], "New name")
	tests.AssertEqual(t, resp["country"], "New country")
}

func TestPostTeam(t *testing.T) {
	setupTest()
	token := getAdminUserToken(t, "test@gmail.com")
	resp, err := doGetRequest("users/me", token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	payload := map[string]interface{}{
		"owner":   int(resp["id"].(float64)),
		"country": "argentina",
		"name":    "los pumas",
		"budget":  100000,
	}

	resp, err = doPostRequest("teams", token, payload, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	teamRes := "teams/" + strconv.Itoa(int(resp["id"].(float64)))
	resp, err = doGetRequest(teamRes, token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTeam(t *testing.T) {
	setupTest()
	token := getAdminUserToken(t, "test@gmail.com")

	resp, err := doGetRequest("users/me", token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	team := resp["team"].(map[string]interface{})
	teamRes := "teams/" + strconv.Itoa(int(team["id"].(float64)))
	/*resp, err = doDeleteRequest(teamRes, token, http.StatusBadRequest)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = doGetRequest(teamRes, token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	players := resp["players"].([]interface{})
	for _, p := range players {
		deletePlayer(t, token, int(p.(float64)))
	}*/

	resp, err = doDeleteRequest(teamRes, token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = doGetRequest(teamRes, token, http.StatusNotFound)
	if err != nil {
		t.Fatal(err)
	}
}
