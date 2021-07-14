package app

import (
	"./models"
	"fmt"
	"gorm.io/gorm/utils/tests"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
)

func TestGetTransfer(t *testing.T) {
	setupTest()
	ask := 10000
	token, playerId, transferId := createTransfer(t, ask)
	resp, err := doGetRequest("transfers/"+strconv.Itoa(transferId), token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	tests.AssertEqual(t, resp["ask"], ask)
	tests.AssertEqual(t, (resp["player"].(map[string]interface{}))["id"], playerId)
}

func TestGetAllTransfers(t *testing.T) {
	setupTest()
	token, players := getTokenAndPlayerIds(t, false)
	asks := make([]int, 0)
	transfers := make([]int, 0)

	for i := 0; i < 10; i++ {
		ask := int(10000 + rand.Float64()*10000)
		transferId := createTransferUsing(t, ask, token, players[i])
		asks = append(asks, ask)
		transfers = append(transfers, transferId)
	}
	resp, err := doGetRequest("transfers", token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	arr := resp["transfers"].([]interface{})
	tests.AssertEqual(t, len(arr), 10)
}

func TestCreateTransfer(t *testing.T) {
	setupTest()
	_, _, _ = createTransfer(t, 10000)
}

func TestCreateTransferOnSamePlayerFails(t *testing.T) {
	setupTest()
	token, players := getTokenAndPlayerIds(t, false)
	_ = createTransferUsing(t, 10, token, players[0])
	_, err := doPostRequest("transfers", token, map[string]interface{}{
		"player_id": players[0],
		"ask":       10000,
	}, http.StatusBadRequest)
	if err != nil {
		t.Fatal(err)
	}
}

func TestBuyTransfer(t *testing.T) {
	setupTest()
	ask := 10000
	token1, playerId, transferId := createTransfer(t, ask)
	token2 := getUserToken(t, "hola@test.com")
	_, err := doPutRequest("transfers/"+strconv.Itoa(transferId)+"/buy", token2, map[string]interface{}{}, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	resp1, err := doGetRequest("users/me/team", token1, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	resp2, err := doGetRequest("users/me/team", token2, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	players1 := resp1["players"].([]interface{})
	players2 := resp2["players"].([]interface{})
	fmt.Println(len(players1))
	fmt.Println(len(players2))
	for _, m := range players1 {
		p := m.(map[string]interface{})
		if int(p["id"].(float64)) == playerId {
			t.Fatal("player is still is the src team")
		}
	}
	isIn := false
	for _, m := range players2 {
		p := m.(map[string]interface{})
		if int(p["id"].(float64)) == playerId {
			isIn = true
			break
		}
	}
	if !isIn {
		t.Fatal("player is not in the dst team")
	}
	// Assert team budgets
	tests.AssertEqual(t, resp1["budget"], models.DefaultTeamBudget+ask)
	tests.AssertEqual(t, resp2["budget"], models.DefaultTeamBudget-ask)

	_, err = doGetRequest("transfers/"+strconv.Itoa(transferId), token1, http.StatusNotFound)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCantBuyTransferIfNotEnoughMoney(t *testing.T) {
	setupTest()
	_, _, transferId := createTransfer(t, 100000000000)
	token2 := getUserToken(t, "hola@test.com")
	_, err := doPutRequest("transfers/"+strconv.Itoa(transferId)+"/buy", token2, map[string]interface{}{}, http.StatusBadRequest)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEditTransfer(t *testing.T) {
	setupTest()
	token, _, transferId := createTransfer(t, 10000)
	_, err := doPatchRequest("transfers/"+strconv.Itoa(transferId), token, map[string]interface{}{
		"ask": 1000000000,
	}, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := doGetRequest("transfers/"+strconv.Itoa(transferId), token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	tests.AssertEqual(t, resp["ask"], 1000000000)
}

func TestEditNotOwnedTransfer(t *testing.T) {
	setupTest()
	_, _, transferId := createTransfer(t, 10000)
	token := getUserToken(t, "hola@tes.com")
	_, err := doPatchRequest("transfers/"+strconv.Itoa(transferId), token, map[string]interface{}{
		"ask": 1000000000,
	}, http.StatusUnauthorized)

	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteTransfer(t *testing.T) {
	setupTest()
	token, _, transferId := createTransfer(t, 10000)

	_, err := doDeleteRequest("transfers/"+strconv.Itoa(transferId), token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	_, err = doGetRequest("transfers/"+strconv.Itoa(transferId), token, http.StatusNotFound)
	if err != nil {
		t.Fatal("found deleted transfer")
	}
}

func createTransfer(t *testing.T, ask int) (string, int, int) {
	token, players := getTokenAndPlayerIds(t, false)
	resp, err := doPostRequest("transfers", token, map[string]interface{}{
		"player_id": players[0],
		"ask":       ask,
	}, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	return token, players[0], int(resp["id"].(float64))
}

func createTransferUsing(t *testing.T, ask int, token string, playerId int) int {
	resp, err := doPostRequest("transfers", token, map[string]interface{}{
		"player_id": playerId,
		"ask":       ask,
	}, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	return int(resp["id"].(float64))
}
