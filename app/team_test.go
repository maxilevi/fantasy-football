package app

import (
	"strconv"
	"testing"
)

func TestQueryUserAndTeamInformation(t *testing.T) {
	token := getUserToken(t, "test@gmail.com")
	resp, err := doGetRequest("user", token)

	if err != nil || resp["email"].(string) != "test@gmail.com" {
		t.Fatal(err)
	}

	teamId := strconv.Itoa(int(resp["team"].(float64)))
	resp, err = doGetRequest("team/"+teamId, token)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { truncateDb() })
}

func TestPatchingTeamInformation(t *testing.T) {
	token := getUserToken(t, "test@gmail.com")
	resp, _ := doGetRequest("user", token)
	res := "team/" + strconv.Itoa(int(resp["team"].(float64)))
	resp, err := doPatchRequest(res, token, map[string]string{
		"name":    "New name",
		"country": "New country",
	})
	if err != nil {
		t.Fatal(err)
	}

	resp, err = doGetRequest(res, token)
	if err != nil || resp["name"] != "New name" || resp["country"] != "New country" {
		t.Fatal(err, resp["name"], resp["country"])
	}

	t.Cleanup(func() { truncateDb() })
}
