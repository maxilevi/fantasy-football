package app

import (
	"net/http"
	"testing"
)

func TestRegisteringUserAndCreatingNewSession(t *testing.T) {
	assertOkRegisteringUser(t, "test@gmail.com", "test1234")
	token := assertOkCreatingSession(t, "test@gmail.com", "test1234")
	if token == "" {
		t.Fatal("invalid token")
	}
	t.Cleanup(func() { truncateDb() })
}

func TestCantLoginWithWrongPassword(t *testing.T) {
	assertOkRegisteringUser(t, "test@gmail.com", "test1234")
	resp, err := doPostRequest("session", map[string]string{
		"email":    "test@gmail.com",
		"password": "asd12345",
	}, http.StatusUnauthorized)
	if err != nil {
		t.Fatal(err)
	}
	val, ok := resp["error"].(bool)
	if !ok || !val {
		t.Fatal("was able to login with incorrect password")
	}
	_, ok = resp["token"].(string)
	if ok {
		t.Fatal("returned a valid token")
	}
	t.Cleanup(func() { truncateDb() })
}

func TestFailRegisteringUser(t *testing.T) {
	assertFailureWhenRegisteringUserWithMessage(t, "test@gmail.com", "12", "Password needs a minimum of at least 8 characters", http.StatusBadRequest)
	assertFailureWhenRegisteringUserWithMessage(t, "test", "12345678", "Invalid email", http.StatusBadRequest)
	t.Cleanup(func() { truncateDb() })
}

func TestCantCreateUserTwice(t *testing.T) {
	assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	assertFailureWhenRegisteringUserWithMessage(t, "test@gmail.com", "dasasdasd2", "Provided email is already registered", http.StatusBadRequest)
	t.Cleanup(func() { truncateDb() })
}

func TestCantQueryUserIfNotLoggedIn(t *testing.T) {
	resp, err := doGetRequest("user", "", http.StatusUnauthorized)
	if err != nil || !resp["error"].(bool) {
		t.Fatal(err)
	}

	t.Cleanup(func() { truncateDb() })
}