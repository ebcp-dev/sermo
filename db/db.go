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
		id uuid DEFAULT uuid_generate_v4 () unique,
		email varchar(225) NOT NULL UNIQUE,
		password varchar(225) NOT NULL,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		primary key (id)
	);
`

// Schema for data table.
const DATA_SCHEMA = `
	CREATE TABLE IF NOT EXISTS data (
		id uuid DEFAULT uuid_generate_v4 () unique,
		strattr varchar(225) NOT NULL UNIQUE,
		intattr int NOT NULL,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		primary key (id)
	);
`

// Receives database credentials and connects to database.
func (db *DB) Initialize(user, password, dbhost, dbname string) {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", user, password, dbhost, dbname)

	var err error
	db.Database, err = sql.Open("postgres", connectionString)
	// Log errors.
	if err != nil {
		log.Fatal(err)
	}
	db.Database.Exec(DB_SETUP)
	db.Database.Exec(USER_SCHEMA)
	db.Database.Exec(DATA_SCHEMA)
}
