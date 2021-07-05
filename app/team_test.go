package app

import (
	"net/http"
	"strconv"
	"testing"
)

func TestQueryUserAndTeamInformation(t *testing.T) {
	setupTest()
	token := getUserToken(t, "test@gmail.com")
	resp, err := doGetRequest("user", token, http.StatusOK)

	if err != nil || resp["email"].(string) != "test@gmail.com" {
		t.Fatal(err)
	}

	teamId := strconv.Itoa(int(resp["team"].(float64)))
	resp, err = doGetRequest("team/"+teamId, token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPatchingTeamInformation(t *testing.T) {
	setupTest()
	token := getUserToken(t, "test@gmail.com")
	resp, _ := doGetRequest("user", token, http.StatusOK)
	res := "team/" + strconv.Itoa(int(resp["team"].(float64)))
	resp, err := doPatchRequest(res, token, map[string]string{
		"name":    "New name",
		"country": "New country",
	}, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	resp, err = doGetRequest(res, token, http.StatusOK)
	if err != nil || resp["name"] != "New name" || resp["country"] != "New country" {
		t.Fatal(err, resp["name"], resp["country"])
	}
}
