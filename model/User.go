package model

import (
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

// Defines user model.
type User struct {
	ID        uuid.UUID `json:"id" sql:"uuid"`
	Email     string    `json:"email" validate:"required" sql:"email"`
	Password  string    `json:"password" validate:"required" sql:"password"`
	CreatedAt time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt time.Time `json:"updatedat" sql:"updatedat"`
	jwt.StandardClaims
}

// Query operations

// Gets a specific user by id.
func (u *User) GetUser(db *sql.DB) error {
	return db.QueryRow("SELECT email, password, createdat, updatedat FROM users WHERE id=$1",
		u.ID).Scan(&u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}

// Gets a specific user by email and password.
func (u *User) GetUserByEmail(db *sql.DB) error {
	return db.QueryRow("SELECT id, email, password, createdat, updatedat FROM users WHERE email=$1 AND password=$2", u.Email, u.Password).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
}

// Gets multiple users. Limit count and start position in db.
func GetUsers(db *sql.DB, start, count int) ([]User, error) {
	rows, err := db.Query(
		"SELECT id, email, password, createdat, updatedat FROM users LIMIT $1 OFFSET $2",
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
		if err := rows.Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}

// CRUD operations

// Create new user and insert to database.
func (u *User) CreateUser(db *sql.DB) error {
	// Scan db after creation if user exists using new user's id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO users(email, password, createdat, updatedat) VALUES($1, $2, $3, $4) RETURNING id, email, password, createdat, updatedat", u.Email, u.Password, timestamp, timestamp).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific user's details by id.
func (u *User) UpdateUser(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE users SET email=$1, password=$2, updatedat=$3 WHERE id=$4 RETURNING id, email, password, createdat, updatedat", u.Email, u.Password, timestamp, u.ID)

	return err
}

// Deletes a specific user by id.
func (u *User) DeleteUser(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM users WHERE id=$1", u.ID)

	return err
}
