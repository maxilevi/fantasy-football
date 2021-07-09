package app

import (
	"./models"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

const testAddr = "localhost:8080"

var app *App

func TestMain(m *testing.M) {
	app = setupTestApp()
	log.SetOutput(ioutil.Discard)
	code := m.Run()
	app.Close()
	os.Exit(code)
}

func getUserToken(t *testing.T, email string) string {
	assertOkRegisteringUser(t, email, "test1234")
	return assertOkCreatingSession(t, email, "test1234")
}

func getAdminUserToken(t *testing.T, email string) string {
	token := getUserToken(t, email)
	res := app.db.Table("users").Where("email = ?", email).Update("permission_level", 1)
	if res.Error != nil {
		t.Fatal(res.Error)
	}
	return token
}

func truncateDb() {
	app.db.Where("1 = 1").Delete(&models.Player{})
	app.db.Where("1 = 1").Delete(&models.Team{})
	app.db.Where("1 = 1").Delete(&models.User{})
	app.db.Where("1 = 1").Delete(&models.Transfer{})
}

func setupTest() {
	truncateDb()
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

func assertOkRegisteringUser(t *testing.T, email string, pass string) int {
	resp, err := doPostRequest("user", "", map[string]interface{}{
		"email":    email,
		"password": pass,
	}, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	if status, valid := resp["code"].(float64); int(status) != 200 || !valid {
		t.Fatalf("unexpected response %v", resp)
	}
	return int(resp["id"].(float64))
}

func assertOkCreatingSession(t *testing.T, email string, pass string) string {
	resp, err := doPostRequest("session", "", map[string]interface{}{
		"email":    email,
		"password": pass,
	}, http.StatusOK)
	if err != nil {
		t.Fatal(err)
	}
	if status, valid := resp["code"].(float64); int(status) != 200 || !valid {
		t.Fatal("unexpected response")
	}
	return resp["token"].(string)
}

func assertFailureWhenRegisteringUserWithMessage(t *testing.T, email string, pass string, msg string, statusCode int) {
	resp, err := doPostRequest("user", "", map[string]interface{}{
		"email":    email,
		"password": pass,
	}, statusCode)
	if err != nil {
		t.Fatal(err)
	}

	if status, valid := resp["code"].(float64); int(status) == 200 || !valid {
		t.Fatal("unexpected response")
	}
	if resp["message"] != msg {
		t.Fatalf("wrong error msg %v", resp["message"])
	}
}

func doPostRequest(resource string, token string, body map[string]interface{}, expectedStatusCode int) (map[string]interface{}, error) {
	return doRequest(resource, token, "POST", body, expectedStatusCode)
}

func doGetRequest(resource string, token string, expectedStatusCode int) (map[string]interface{}, error) {
	return doRequest(resource, token, "GET", map[string]interface{}{}, expectedStatusCode)
}

func doPatchRequest(resource string, token string, body map[string]interface{}, expectedStatusCode int) (map[string]interface{}, error) {
	return doRequest(resource, token, "PATCH", body, expectedStatusCode)
}

func doDeleteRequest(resource string, token string, expectedStatusCode int) (map[string]interface{}, error) {
	return doRequest(resource, token, "DELETE", map[string]interface{}{}, expectedStatusCode)
}

func doRequest(resource string, token string, method string, body map[string]interface{}, expectedStatusCode int) (map[string]interface{}, error) {

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
		req.Header.Add("Content-Type", "application/json")
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != expectedStatusCode {
		return nil, fmt.Errorf("unexpected status code %v", resp.StatusCode)
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
