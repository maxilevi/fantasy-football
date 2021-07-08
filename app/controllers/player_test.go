package controllers

import (
	"bytes"
	"encoding/json"
	"gorm.io/gorm/utils/tests"
	"net/http"
	"net/http/httptest"
	"testing"
	"../models"
)

func TestGetPlayerJsonFromRequestDoesntOverrideTeamId(t *testing.T) {
	c := &PlayerController{}
	body := map[string]interface{}{
		"age": 10,
	}
	jsonBody, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "", bytes.NewBuffer(jsonBody))

	payload, err := c.getPlayerJsonFromRequest(w, req)
	if err != nil {
		t.Fatal(err)
	}

	player := models.Player{}
	player.TeamID = 1
	c.fillPlayerData(&player, payload, true)

	tests.AssertEqual(t, player.TeamID, 1)
	tests.AssertEqual(t, player.Age, 10)
}

func TestFillPlayerDataDoesNotChangeIfNotAdmin(t *testing.T) {
	c := &PlayerController{}
	player := models.Player{}

	c.fillPlayerData(&player, playerPayload{
		Age: 10,
		Position: 123,
		MarketValue: 101000,
		Team: 100,
	}, false)

	tests.AssertEqual(t, player.TeamID, 0)
	tests.AssertEqual(t, player.MarketValue, 0)
	tests.AssertEqual(t, player.Position, 0)
	tests.AssertEqual(t, player.Age, 0)
}

func TestFillPlayerDataChangesIfAdmin(t *testing.T) {
	c := &PlayerController{}
	player := models.Player{}

	c.fillPlayerData(&player, playerPayload{
		Age: 10,
		Position: 123,
		MarketValue: 101000,
		Team: 100,
	}, true)

	tests.AssertEqual(t, player.TeamID, 10)
	tests.AssertEqual(t, player.MarketValue, 123)
	tests.AssertEqual(t, player.Position, 101000)
	tests.AssertEqual(t, player.Age, 100)
}

func TestFillPlayerFillsData(t *testing.T) {
	c := &PlayerController{}
	player := models.Player{}
	player.TeamID = 1

	c.fillPlayerData(&player, playerPayload{}, false)

	tests.AssertEqual(t, player.TeamID, 1)
}