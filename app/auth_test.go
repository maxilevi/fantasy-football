package app

import (
	"net/http"
	"testing"
)

func TestRegisteringUserAndCreatingNewSession(t *testing.T) {
	setupTest()
	assertOkRegisteringUser(t, "test@gmail.com", "test1234")
	token := assertOkCreatingSession(t, "test@gmail.com", "test1234")
	if token == "" {
		t.Fatal("invalid token")
	}
}

func TestCantLoginWithWrongPassword(t *testing.T) {
	setupTest()
	assertOkRegisteringUser(t, "test@gmail.com", "test1234")
	resp, err := doPostRequest("session", "", map[string]interface{}{
		"email":    "test@gmail.com",
		"password": "asd12345",
	}, http.StatusUnauthorized)
	if err != nil {
		t.Fatal(err)
	}
	val, ok := resp["status"].(float64)
	if !ok || int(val) == 200 {
		t.Fatal("was able to login with incorrect password")
	}
	_, ok = resp["token"].(string)
	if ok {
		t.Fatal("returned a valid token")
	}
}

func TestFailRegisteringUser(t *testing.T) {
	setupTest()
	assertFailureWhenRegisteringUserWithMessage(t, "test@gmail.com", "12", "Password needs a minimum of at least 8 characters", http.StatusBadRequest)
	assertFailureWhenRegisteringUserWithMessage(t, "test", "12345678", "Invalid email", http.StatusBadRequest)
}

func TestCantCreateUserTwice(t *testing.T) {
	setupTest()
	assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	assertFailureWhenRegisteringUserWithMessage(t, "test@gmail.com", "dasasdasd2", "Provided email is already registered", http.StatusBadRequest)
}

func TestCantQueryUserIfNotLoggedIn(t *testing.T) {
	setupTest()
	resp, err := doGetRequest("user", "", http.StatusUnauthorized)
	if err != nil || int(resp["status"].(float64)) == 200 {
		t.Fatal(err)
	}
}
