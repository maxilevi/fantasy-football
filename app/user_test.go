package app

import (
	"gorm.io/gorm/utils/tests"
	"net/http"
	"strconv"
	"testing"
)

func TestGetMe(t *testing.T) {
	setupTest()
	email := "test@gmail.com"
	token := getUserToken(t, email)
	resp, err := doGetRequest("users/me", token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	tests.AssertEqual(t, resp["email"].(string), email)
}

func TestGetUser(t *testing.T) {
	setupTest()
	userId := assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	token := getAdminUserToken(t, "admin@gmail.com")
	resp, err := doGetRequest("users/"+strconv.Itoa(userId), token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	tests.AssertEqual(t, resp["email"].(string), "test@gmail.com")
}

func TestCantGetUserNotAdmin(t *testing.T) {
	setupTest()
	userId := assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	token := getUserToken(t, "admin@gmail.com")

	_, err := doGetRequest("users/"+strconv.Itoa(userId), token, http.StatusUnauthorized)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteUser(t *testing.T) {
	setupTest()
	userId := assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	token := getAdminUserToken(t, "admin@gmail.com")

	_, err := doDeleteRequest("users/"+strconv.Itoa(userId), token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	_, err = doGetRequest("users/"+strconv.Itoa(userId), token, http.StatusNotFound)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPatchUser(t *testing.T) {
	setupTest()
	userId := assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	token := getAdminUserToken(t, "admin@gmail.com")

	payload := map[string]interface{}{
		"email": "new_test@gmail.com",
	}

	_, err := doPatchRequest("users/"+strconv.Itoa(userId), token, payload, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := doGetRequest("users/"+strconv.Itoa(userId), token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	for key, value := range resp {
		if key == "code" || key == "team" || key == "id" {
			continue
		}
		tests.AssertEqual(t, value, payload[key])
	}
}
