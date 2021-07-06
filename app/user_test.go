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
	resp, err := doGetRequest("user", token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	tests.AssertEqual(t, resp["email"].(string), email)
}

func TestGetUser(t *testing.T) {
	setupTest()
	userId := assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	token := getAdminUserToken(t, "admin@gmail.com")
	resp, err := doGetRequest("user/" + strconv.Itoa(userId), token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	tests.AssertEqual(t, resp["email"].(string), "test@gmail.com")
}

func TestCantGetUserNotAdmin(t *testing.T) {
	setupTest()
	userId := assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	token := getUserToken(t, "admin@gmail.com")

	_, _ = doGetRequest("user/" + strconv.Itoa(userId), token, http.StatusUnauthorized)
}

func TestDeleteUser(t *testing.T) {
	setupTest()
	userId := assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	token := getAdminUserToken(t, "admin@gmail.com")

	_, _ = doDeleteRequest("user/"+strconv.Itoa(userId), token, http.StatusOK)
	_, _ = doGetRequest("user/"+strconv.Itoa(userId), token, http.StatusNotFound)
}

func TestPatchUser(t *testing.T) {
	setupTest()
	userId := assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	token := getAdminUserToken(t, "admin@gmail.com")

	payload := map[string]interface{}{
		"email": "new_test@gmail.com",
	}

	_, err := doPatchRequest("user/"+strconv.Itoa(userId), token, payload, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := doGetRequest("user/"+strconv.Itoa(userId), token, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}

	for key, value := range resp {
		if key == "error" {
			continue
		}
		tests.AssertEqual(t, value, payload[key])
	}
}