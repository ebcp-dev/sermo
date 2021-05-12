package db

import (
	"database/sql"
	"fmt"
	"log"
)

type DB struct {
	Database *sql.DB
}

// Schema for user table
const DB_SCHEMA = `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	CREATE TABLE IF NOT EXISTS users (
		id uuid DEFAULT uuid_generate_v4 () unique,
		email varchar(225) NOT NULL UNIQUE,
		password varchar(225) NOT NULL,
		createdat timestamp NOT NULL,
		updatedat timestamp NOT NULL,
		primary key (id)
	);
`

// Receives database credentials and connects to database.
func (db *DB) Initialize(user, password, dbname string) {
	connectionString :=
		fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	db.Database, err = sql.Open("postgres", connectionString)
	// Log errors.
	if err != nil {
		log.Fatal(err)
	}
	db.Database.Exec(DB_SCHEMA)
}
