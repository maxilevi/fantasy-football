package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
)

const testAddr = "localhost:8080"

func TestRegisteringUserAndCreatingNewSession(t *testing.T) {
	app := setupTestApp(t)
	defer app.Close()

	assertOkRegisteringUser(t, "test@gmail.com", "test1234")
	token := assertOkCreatingSession(t, "test@gmail.com", "test1234")
	fmt.Println(token)
}

func TestFailRegisteringUser(t *testing.T) {
	app := setupTestApp(t)
	defer app.Close()

	assertFailureWhenRegisteringUserWithMessage(t, "test@gmail.com", "12", "Password needs a minimum of at least 8 characters")
	assertFailureWhenRegisteringUserWithMessage(t, "test", "12345678", "Invalid email")
	assertFailureWhenRegisteringUserWithMessage(t, "test@gmail.com", "12", "Incorrect body parameters")
}

func TestCantCreateUserTwice(t *testing.T) {
	app := setupTestApp(t)
	defer app.Close()

	assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	assertFailureWhenRegisteringUserWithMessage(t, "test@gmail.com", "32413", "Provided email is already registered")
}

func setupTestApp(t *testing.T) *App {
	app, err := CreateApp(testAddr, os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PORT"))
	if err != nil {
		t.Error(err)
	}
	go app.Run()
	return app
}

func assertOkRegisteringUser(t *testing.T, email string, pass string) {
	resp, err := doPostRequest("user", map[string]string{
		"email": email,
		"password": pass,
	})
	if err != nil {
		t.Error(err)
	}
	if ok, err := strconv.ParseBool(resp["error"]); !ok || err != nil {
		t.Error("unexpected response")
	}
}

func assertOkCreatingSession(t *testing.T, email string, pass string) string {
	resp, err := doPostRequest("session", map[string]string{
		"email": email,
		"password": pass,
	})
	if err != nil {
		t.Error(err)
	}
	if ok, err := strconv.ParseBool(resp["error"]); !ok || err != nil {
		t.Error("unexpected response")
	}
	return resp["token"]
}

func assertFailureWhenRegisteringUserWithMessage(t *testing.T, email string, pass string, msg string) {
	resp, err := doPostRequest("user", map[string]string{
		"email": email,
		"password": pass,
	})
	if err != nil {
		t.Error(err)
	}
	if ok, err := strconv.ParseBool(resp["error"]); ok || err != nil {
		t.Error("unexpected response")
	}
	if resp["message"] != msg {
		t.Error("wrong error msg")
	}
}

func doPostRequest(resource string, body map[string]string) (map[string]string, error) {
	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	responseBody := bytes.NewBuffer(postBody)
	resp, err := http.Post("http://" + testAddr + "/api/" + resource, "application/json", responseBody)
	if err != nil {
		return nil, err
	}
	bodystr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var m map[string]string
	err = json.Unmarshal(bodystr, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}