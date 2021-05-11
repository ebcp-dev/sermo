package test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/ebcp-dev/gorest-api/app"
	"github.com/google/uuid"
)

// References App struct in app.go.
var a app.App

//Generate new uuid for test
var testID = uuid.NewString()

// Executes before all other tests.
func TestMain(m *testing.M) {
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"))
	ensureTableExists()
	// Executes tests.
	code := m.Run()
	// Cleans testing table is cleaned from database.
	clearTable()
	os.Exit(code)
}

// Test functions

// Tests response if users table is empty.
// Deletes all records from users table and sends GET request to /users endpoint.
func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

// Test response if requested user is non-existent.
// Tests if status code = 404 & response message = "User not found".
func TestGetNonExistentUser(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/user/"+testID, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "User not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'User not found'. Got '%s'", m["error"])
	}
}

// Test response on login route.
// Tests if status code = 200.
func TestLoginUser(t *testing.T) {
	clearTable()
	addUsers(1)

	var jsonStr = []byte(`{"email":"testemail1@gmail.com", "password":"password1"}`)
	req, _ := http.NewRequest("POST", "/user/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test response when fetching a specific user.
// Tests if status code = 200.
func TestGetUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/"+testID, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test the process of creating a new user by manually adding a test user to db.
// Tests if status code = 200 & response contains JSON object with the right contents.
func TestCreateUser(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"email":"testemail1@gmail.com", "password": "password1"}`)
	req, _ := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["email"] != "testemail1@gmail.com" {
		t.Errorf("Expected user email to be 'testemail1@gmail.com'. Got '%v'", m["email"])
	}

	if m["password"] != "password1" {
		t.Errorf("Expected user password to be 'password1'. Got '%v'", m["password"])
	}
}

// Test process of updating a user.
// Tests if status code = 200 & response contains JSON object with the updated contents.
func TestUpdateUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/"+testID, nil)
	response := executeRequest(req)
	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)

	var jsonStr = []byte(`{"email":"testemail1@gmail.com - updated email", "password": "password1 - updated password"}`)
	req, _ = http.NewRequest("PUT", "/user/"+testID, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalUser["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
	}

	if m["email"] == originalUser["email"] {
		t.Errorf("Expected the email to change from '%v' to '%v'. Got '%v'", originalUser["email"], m["email"], m["email"])
	}

	if m["password"] == originalUser["password"] {
		t.Errorf("Expected the password to change from '%v' to '%v'. Got '%v'", originalUser["password"], m["password"], m["password"])
	}
}

// Test process of deleting users.
// Tests if status code = 200.
func TestDeleteUser(t *testing.T) {
	clearTable()
	addUsers(1)

	req, _ := http.NewRequest("GET", "/user/"+testID, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/user/"+testID, nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/user/"+uuid.NewString(), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// Helper functions

// Executes http request using the router and returns response.
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

// Compares actual response to expected response of tested http request.
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Adds 1 or more records to table for testing.
func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 1; i <= count; i++ {
		timestamp := time.Now()
		a.DB.Exec("INSERT INTO users(id, email, password, createdat, updatedat) VALUES($1, $2, $3, $4, $5)", testID, "testemail"+strconv.Itoa(i)+"@gmail.com", "password"+strconv.Itoa(i), timestamp, timestamp)
	}
}

// Ensures table needed for testing exists.
func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

// Clean test table.
func clearTable() {
	a.DB.Exec("DELETE FROM users")
}

// SQL query to create table.
const tableCreationQuery = `
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT uuid_generate_v4 () unique,
		email varchar(225) NOT NULL UNIQUE,
		password varchar(225) NOT NULL,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		primary key (id)
	);`
