package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ebcp-dev/sermo/app/auth"
)

// Test functions

// Tests response if users table is empty.
// Deletes all records from users table and sends GET request to /users endpoint.
func TestEmptyUserTable(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := auth.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	req, _ := http.NewRequest("GET", "/api/users", nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
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
	// Generate JWT for authorization.
	validToken, err := auth.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/api/user/"+userTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
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
	req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test response when fetching a specific user.
// Tests if status code = 200.
func TestGetUser(t *testing.T) {
	clearTable()
	addUsers(1)
	// Generate JWT for authorization.
	validToken, err := auth.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/api/user/"+userTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test the process of creating a new user by manually adding a test user to db.
// Tests if status code = 200 & response contains JSON object with the right contents.
func TestCreateUser(t *testing.T) {
	clearTable()

	var jsonStr = []byte(`{"email":"testemail1@gmail.com", "password": "password1"}`)
	req, _ := http.NewRequest("POST", "/api/user", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["email"] != "testemail1@gmail.com" {
		t.Errorf("Expected user email to be 'testemail1@gmail.com'. Got '%v'", m["email"])
	}
	// Get string value from map index.
	passwordHash := m["password"].(string)
	if !auth.ComparePasswords(passwordHash, []byte("password1")) {
		t.Errorf("Passwords do not match.")
	}
}

// Test process of updating a user.
// Tests if status code = 200 & response contains JSON object with the updated contents.
func TestUpdateUser(t *testing.T) {
	clearTable()
	addUsers(1)
	// Generate JWT for authorization.
	validToken, err := auth.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/api/user/"+userTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)

	var jsonStr = []byte(`{"email":"testemail1@gmail.com - updated email", "password": "password1 - updated password"}`)
	req, _ = http.NewRequest("PUT", "/api/user/"+userTestID.String(), bytes.NewBuffer(jsonStr))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["userid"] != originalUser["userid"] {
		t.Errorf("Expected the userid to remain the same (%v). Got %v", originalUser["userid"], m["userid"])
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
	// Generate JWT for authorization.
	validToken, err := auth.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	// Check that user exists.
	req, _ := http.NewRequest("GET", "/api/user/"+userTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Delete user.
	req, _ = http.NewRequest("DELETE", "/api/user/"+userTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Check if user still exists.
	req, _ = http.NewRequest("GET", "/api/user/"+userTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// Helper functions

// Adds 1 or more records to table for testing.
func addUsers(count int) {
	if count < 1 {
		count = 1
	}

	for i := 1; i <= count; i++ {
		timestamp := time.Now()
		passwordHash := auth.HashAndSalt([]byte("password" + strconv.Itoa(i)))
		d.Database.Exec("INSERT INTO users(userid, email, password, createdat, updatedat) VALUES($1, $2, $3, $4, $5)", userTestID, "testemail"+strconv.Itoa(i)+"@gmail.com", passwordHash, timestamp, timestamp)
	}
}
