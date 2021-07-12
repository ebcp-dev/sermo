package db

import (
	"database/sql"
	"fmt"
	"log"
)

type DB struct {
	Database *sql.DB
}

const DB_SETUP = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
`

// Schema for user table.
const USER_SCHEMA = `
	CREATE TABLE IF NOT EXISTS users (
		userid UUID DEFAULT uuid_generate_v4 () UNIQUE,
		email VARCHAR(90) NOT NULL UNIQUE,
		password VARCHAR(100) NOT NULL,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		PRIMARY KEY (userid)
	);
`

// Schema for data table.
const CHANNEL_SCHEMA = `
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

// Receives database credentials and connects to database.
func (db *DB) Initialize(user string, password string, dbhost string, dbname string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, password, dbhost, dbname)

	var err error
	db.Database, err = sql.Open("postgres", connectionString)
	// Log errors.
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connecting to db '%v' as user '%v'.", dbname, user)
	db.Database.Exec(DB_SETUP)
	db.Database.Exec(USER_SCHEMA)
	db.Database.Exec(CHANNEL_SCHEMA)
}
