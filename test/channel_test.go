package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ebcp-dev/sermo/app/utils"
	model "github.com/ebcp-dev/sermo/models"
)

// Test functions

// Tests response if channel table is empty.
// Deletes all records from channel table and sends GET request to /channel endpoint.
func TestEmptyChannelTable(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := utils.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	req, _ := http.NewRequest("GET", "/channels", nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

// Test response if requested channel is non-existent.
// Tests if status code = 404 & response message = "Channel not found".
func TestGetNonExistentChannel(t *testing.T) {
	clearTable()
	// Generate JWT for authorization.
	validToken, err := utils.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/channel/"+channelTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Channel not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Channel not found'. Got '%s'", m["error"])
	}
}

// Test response when fetching a specific channel.
// Tests if status code = 200.
func TestGetChannel(t *testing.T) {
	clearTable()
	addChannel(1)
	// Generate JWT for authorization.
	validToken, err := utils.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	req, _ := http.NewRequest("GET", "/channel/"+channelTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// Test the process of creating a new channel by manually adding a test channel to db.
// Tests if status code = 200 & response contains JSON object with the right contents.
func TestCreateChannel(t *testing.T) {
	clearTable()
	// Create new user for foreign key constraint.
	addUsers(1)
	// Generate JWT for authorization.
	validToken, err := utils.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}

	// var jsonStr = []byte(`{"channelname":"channel1", "maxpopulation": 1`)
	newChannel := model.Channel{
		ChannelName:   "channel1",
		MaxPopulation: 1,
		UserID:        userTestID,
	}
	payload, err := json.Marshal(newChannel)
	if err != nil {
		t.Error("Failed to parse JSON")
	}
	req, _ := http.NewRequest("POST", "/channel", bytes.NewBuffer(payload))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["channelname"] != "channel1" {
		t.Log(m["channelname"])
		t.Errorf("Expected channel channelname to be 'channel1'. Got '%v'", m["channelname"])
	}
	// Convert 1 to float64 because Go maps convert int values to float64.
	if m["maxpopulation"] != float64(1) {
		t.Errorf("Expected channel maxpopulation to be 1. Got '%v'", m["maxpopulation"])
	}
}

// Test process of updating a channel.
// Tests if status code = 200 & response contains JSON object with the updated contents.
func TestUpdateChannel(t *testing.T) {
	clearTable()
	addChannel(1)
	// Generate JWT for authorization.
	validToken, err := utils.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	req, _ := http.NewRequest("GET", "/channel/"+channelTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	var originalChannel map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalChannel)

	var jsonStr = []byte(`{"channelname":"channel1 - updated", "maxpopulation": 2}`)
	req, _ = http.NewRequest("PUT", "/channel/"+channelTestID.String(), bytes.NewBuffer(jsonStr))
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	req.Header.Set("Content-Type", "application/json")

	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["channelid"] != originalChannel["channelid"] {
		t.Errorf("Expected the channelid to remain the same (%v). Got %v", originalChannel["channelid"], m["channelid"])
	}

	if m["channelname"] == originalChannel["channelname"] {
		t.Errorf("Expected the channelname to change from '%v' to '%v'. Got '%v'", originalChannel["channelname"], m["channelname"], m["channelname"])
	}

	if m["maxpopulation"] == originalChannel["maxpopulation"] {
		t.Errorf("Expected the maxpopulation to change from '%v' to '%v'. Got '%v'", originalChannel["maxpopulation"], m["maxpopulation"], m["maxpopulation"])
	}
}

// Test process of deleting channel.
// Tests if status code = 200.
func TestDeleteChannel(t *testing.T) {
	clearTable()
	addChannel(1)
	// Generate JWT for authorization.
	validToken, err := utils.GenerateJWT()
	if err != nil {
		t.Error("Failed to generate token")
	}
	// Check that channel exists.
	req, _ := http.NewRequest("GET", "/channel/"+channelTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Delete channel.
	req, _ = http.NewRequest("DELETE", "/channel/"+channelTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	// Check if channel still exists.
	req, _ = http.NewRequest("GET", "/channel/"+channelTestID.String(), nil)
	// Add "Token" header to request with generated token.
	req.Header.Add("Token", validToken)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

// Helper functions

// Adds 1 or more records to table for testing.
func addChannel(count int) {
	timestamp := time.Now()
	if count < 1 {
		count = 1
	}

	// Create new user for foreign key constraint.
	addUsers(1)

	for i := 1; i <= count; i++ {
		d.Database.Exec("INSERT INTO channels(channelid, channelname, maxpopulation, userid, createdat, updatedat) VALUES($1, $2, $3, $4, $5, $6)", channelTestID, "channel"+strconv.Itoa(i), i, userTestID, timestamp, timestamp)
	}
}
