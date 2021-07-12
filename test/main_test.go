package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ebcp-dev/sermo/app/api"
	"github.com/ebcp-dev/sermo/db"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

//Generate new uuid for test
var userTestID = uuid.New()
var channelTestID = uuid.New()

// References App struct in app.go.
var a api.Api

// References DB struct in app.go.
var d db.DB

// Executes before all other tests.
func TestMain(m *testing.M) {
	os.Setenv("ENV", "test")
	// Find and read the config file.
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	db_user := viper.GetString("TEST_DB_USERNAME")
	db_pass := viper.GetString("TEST_DB_PASSWORD")
	db_host := viper.GetString("TEST_DB_HOST")
	db_name := viper.GetString("TEST_DB_NAME")
	a.InitializeAPI()
	d.Initialize(db_user, db_pass, db_host, db_name)

	ensureTableExists()
	// Executes tests.
	code := m.Run()
	// Cleans testing table is cleaned from database.
	clearTable()
	os.Exit(code)
}

// Helpers

// SQL query to create table.
const tableCreationQuery = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TABLE IF NOT EXISTS users (
		userid UUID DEFAULT uuid_generate_v4 () UNIQUE,
		email VARCHAR(90) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		PRIMARY KEY (userid)
	);
	CREATE TABLE IF NOT EXISTS channels (
		channelid UUID DEFAULT uuid_generate_v4 () UNIQUE,
		channelname VARCHAR(20) NOT NULL UNIQUE,
		maxpopulation int NOT NULL DEFAULT 1,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		userid UUID NOT NULL,
		PRIMARY KEY (channelid),
		CONSTRAINT fk_user FOREIGN KEY (userid) 
			REFERENCES users(userid) ON DELETE CASCADE
	);
`

// Executes http request using the router and returns response.
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

// Compares actual response to expected response of http request.
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Ensures table needed for testing exists.
func ensureTableExists() {
	if _, err := d.Database.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

// Clean test tables.
func clearTable() {
	d.Database.Exec("DELETE FROM channels")
	d.Database.Exec("DELETE FROM users")
}
