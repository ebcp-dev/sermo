package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ebcp-dev/gorest-api/api/app"
	"github.com/ebcp-dev/gorest-api/api/model"
	"github.com/google/uuid"
)

// Test functions

// Tests response if data table is empty.
// Deletes all records from data table and sends GET request to /data endpoint.
func TestEmptyDataTable(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	req, _ := http.NewRequest("GET", "/data", nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

// Test response if requested data is non-existent.
// Tests if status code = 404 & response message = "Data not found".
func TestGetNonExistentData(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/data/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Data not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Data not found'. Got '%s'", m["error"])
	}
}

// Test response when fetching a specific data.
// Tests if status code = 200.
func TestGetData(t *testing.T) {
	clearTable()
	addData(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/data/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test the process of creating a new data by manually adding a test data to db.
// Tests if status code = 200 & response contains JSON object with the right contents.
func TestCreateData(t *testing.T) {
	clearTable()

	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	// var jsonStr = []byte(`{"strattr":"string1", "intattr": 1`)
	newData := model.Data{
		StrAttr: "string1",
		IntAttr: 1,
	}
	payload, err := json.Marshal(newData)
	if err != nil {
		t.Error("Failed to parse JSON")
	}
	req, _ := http.NewRequest("POST", "/data", bytes.NewBuffer(payload))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["strattr"] != "string1" {
		t.Errorf("Expected data strattr to be 'string1'. Got '%v'", m["strattr"])
	}
	// Convert 1 to float64 because Go maps convert int values to float64.
	if m["intattr"] != float64(1) {
		t.Errorf("Expected data intattr to be 1. Got '%v'", m["intattr"])
	}
}

// Test process of updating a data.
// Tests if status code = 200 & response contains JSON object with the updated contents.
func TestUpdateData(t *testing.T) {
	clearTable()
	addData(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/data/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	var originalData map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalData)

	var jsonStr = []byte(`{"strattr":"string1 - updated strattr", "intattr": 2}`)
	req, _ = http.NewRequest("PUT", "/data/"+testID, bytes.NewBuffer(jsonStr))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalData["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalData["id"], m["id"])
	}

	if m["strattr"] == originalData["strattr"] {
		t.Errorf("Expected the strattr to change from '%v' to '%v'. Got '%v'", originalData["strattr"], m["strattr"], m["strattr"])
	}

	if m["intattr"] == originalData["intattr"] {
		t.Errorf("Expected the intattr to change from '%v' to '%v'. Got '%v'", originalData["intattr"], m["intattr"], m["intattr"])
	}
}

// Test process of deleting data.
// Tests if status code = 200.
func TestDeleteData(t *testing.T) {
	clearTable()
	addData(1)
	// Generate JWT for authorization.
	validToken, err := app.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	// Check that data exists.
	req, _ := http.NewRequest("GET", "/data/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Delete data.
	req, _ = http.NewRequest("DELETE", "/data/"+testID, nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Check if data still exists.
	req, _ = http.NewRequest("GET", "/data/"+uuid.NewString(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// Helper functions

// Adds 1 or more records to table for testing.
func addData(count int) {
	if count < 1 {
		count = 1
	}

	for i := 1; i <= count; i++ {
		timestamp := time.Now()
		d.Database.Exec("INSERT INTO data(id, strattr, intattr, createdat, updatedat) VALUES($1, $2, $3, $4, $5)", testID, "string"+strconv.Itoa(i), i, timestamp, timestamp)
	}
}
