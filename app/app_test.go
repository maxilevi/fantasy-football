package app

import (
	"./models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"
)

const testAddr = "localhost:8080"

var app *App

func TestMain(m *testing.M) {
	app = setupTestApp()
	truncateDb()
	code := m.Run()
	app.Close()
	os.Exit(code)
}

/* Auth tests */

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
	})
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
	assertFailureWhenRegisteringUserWithMessage(t, "test@gmail.com", "12", "Password needs a minimum of at least 8 characters")
	assertFailureWhenRegisteringUserWithMessage(t, "test", "12345678", "Invalid email")
	t.Cleanup(func() { truncateDb() })
}

func TestCantCreateUserTwice(t *testing.T) {
	assertOkRegisteringUser(t, "test@gmail.com", "12345678")
	assertFailureWhenRegisteringUserWithMessage(t, "test@gmail.com", "dasasdasd2", "Provided email is already registered")
	t.Cleanup(func() { truncateDb() })
}

func TestCantQueryUserIfNotLoggedIn(t *testing.T) {
	resp, err := doGetRequest("user", "")
	if err != nil || !resp["error"].(bool) {
		t.Fatal(err)
	}

	t.Cleanup(func() { truncateDb() })
}

/* Team tests */

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

/* Player tests */

func TestGetPlayerInformation(t *testing.T) {
	token := getUserToken(t, "test@gmail.com")
	resp, err := doGetRequest("player", token)

	if err != nil {
		t.Fatal(err)
	}

	teamId := strconv.Itoa(int(resp["team"].(float64)))
	resp, err = doGetRequest("team/"+teamId, token)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { truncateDb() })
}

func getUserToken(t *testing.T, email string) string {
	assertOkRegisteringUser(t, email, "test1234")
	return assertOkCreatingSession(t, email, "test1234")
}

func truncateDb() {
	app.db.Where("1 = 1").Delete(&models.Player{})
	app.db.Where("1 = 1").Delete(&models.Team{})
	app.db.Where("1 = 1").Delete(&models.User{})
}

func setupTestApp() *App {
	err := godotenv.Load("../.env")
	if err != nil {
		panic("error loading .env file")
	}
	app, err := CreateApp(testAddr, os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PORT"))
	if err != nil {
		panic(err)
	}
	go app.Run()
	timeout := time.Now().Add(5 * time.Second)
	for {
		if app.IsRunning || time.Now().After(timeout) {
			if !app.IsRunning {
				panic("failed to start app (timeout)")
			}
			break
		}
	}
	return app
}

func assertOkRegisteringUser(t *testing.T, email string, pass string) {
	resp, err := doPostRequest("user", map[string]string{
		"email":    email,
		"password": pass,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err, valid := resp["error"].(bool); err || !valid {
		t.Fatalf("unexpected response %v", resp["message"])
	}
}

func assertOkCreatingSession(t *testing.T, email string, pass string) string {
	resp, err := doPostRequest("session", map[string]string{
		"email":    email,
		"password": pass,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err, valid := resp["error"].(bool); err || !valid {
		t.Fatal("unexpected response")
	}
	return resp["token"].(string)
}

func assertFailureWhenRegisteringUserWithMessage(t *testing.T, email string, pass string, msg string) {
	resp, err := doPostRequest("user", map[string]string{
		"email":    email,
		"password": pass,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err, valid := resp["error"].(bool); !err || !valid {
		t.Fatal("unexpected response")
	}
	if resp["message"] != msg {
		t.Fatalf("wrong error msg %v", resp["message"])
	}
}

func doPostRequest(resource string, body interface{}) (map[string]interface{}, error) {
	return doRequest(resource, "", "POST", body)
}

func doGetRequest(resource string, token string) (map[string]interface{}, error) {
	return doRequest(resource, token, "GET", map[string]string{})
}

func doPatchRequest(resource string, token string, body interface{}) (map[string]interface{}, error) {
	return doRequest(resource, token, "PATCH", body)
}

func doRequest(resource string, token string, method string, body interface{}) (map[string]interface{}, error) {

	postBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	requestBody := bytes.NewBuffer(postBody)
	req, err := http.NewRequest(method, "http://"+testAddr+"/api/"+resource, requestBody)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodystr, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(bodystr)
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(bodystr, &m)
	if err != nil {
		fmt.Println(string(bodystr))
		return nil, err
	}

	return m, nil
}
