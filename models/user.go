package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines user model.
type User struct {
	UserID    uuid.UUID `json:"userid" sql:"uuid"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	CreatedAt time.Time `json:"createdat" validate:"required"`
	UpdatedAt time.Time `json:"updatedat" validate:"required"`
}

// Query operations

// Gets a specific user by UserID.
func (u *User) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT email, password, createdat, updatedat FROM users WHERE UserID=$1",
		u.UserID).Scan(&u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}

// Gets a specific user by Email.
func (u *User) GetUserByEmail(db *sql.DB) error {
	return db.QueryRow("SELECT email, password, createdat, updatedat FROM users WHERE email=$1",
		u.Email).Scan(&u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}

// Gets a specific user by email and password.
func (u *User) GetUserByEmailAndPassword(db *sql.DB) error {
	return db.QueryRow("SELECT UserID, email, password, createdat, updatedat FROM users WHERE email=$1 AND password=$2", u.Email, u.Password).Scan(&u.UserID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}

// Gets multiple users. Limit count and start position in db.
func GetUsers(db *sql.DB, start, count int) ([]User, error) {
	rows, err := db.Query(
		"SELECT UserID, email, password, createdat, updatedat FROM users LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	users := []User{}

	// Store query results into users variable if no errors.
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.UserID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// CRUD operations

// Create new user and insert to database.
func (u *User) CreateUser(db *sql.DB) error {
	// Scan db after creation if user exists using new user's UserID.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO users(email, password, createdat, updatedat) VALUES($1, $2, $3, $4) RETURNING UserID, email, password, createdat, updatedat", u.Email, u.Password, timestamp, timestamp).Scan(&u.UserID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific user's details by UserID.
func (u *User) UpdateUser(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE users SET email=$1, password=$2, updatedat=$3 WHERE UserID=$4 RETURNING UserID, email, password, createdat, updatedat", u.Email, u.Password, timestamp, u.UserID)

	return err
}

// Deletes a specific user by UserID.
func (u *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE UserID=$1", u.UserID)

	return err
}
