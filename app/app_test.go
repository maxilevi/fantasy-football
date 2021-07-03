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
