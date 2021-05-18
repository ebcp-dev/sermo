package test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ebcp-dev/sermo/app"
	"github.com/ebcp-dev/sermo/db"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

// References App struct in app.go.
var a app.App

// References DB struct in app.go.
var d db.DB

//Generate new uuid for test
var testID = uuid.NewString()

// Executes before all other tests.
func TestMain(m *testing.M) {
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	// Find and read the config file
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	db_user := viper.GetString("TEST_DB_USERNAME")
	db_pass := viper.GetString("TEST_DB_PASSWORD")
	db_host := viper.GetString("TEST_DB_HOST")
	db_name := viper.GetString("TEST_DB_NAME")

	d.Initialize(db_user, db_pass, db_host, db_name)
	a.Initialize()
	ensureTableExists()
	// Executes tests.
	code := m.Run()
	// Cleans testing table is cleaned from database.
	clearTable()
	os.Exit(code)
}

// Helpers

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
	d.Database.Exec("DELETE FROM users")
	d.Database.Exec("DELETE FROM data")
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
	);
	CREATE TABLE IF NOT EXISTS data (
		id uuid DEFAULT uuid_generate_v4 () unique,
		strattr varchar(225) NOT NULL UNIQUE,
		intattr int NOT NULL,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		primary key (id)
	);
`
