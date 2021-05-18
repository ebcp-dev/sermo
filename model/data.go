package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines data model.
type Data struct {
	ID        uuid.UUID `json:"id" sql:"uuid"`
	StrAttr   string    `json:"strattr" validate:"required" sql:"strattr"`
	IntAttr   int       `json:"intattr" validate:"required" sql:"intattr"`
	CreatedAt time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt time.Time `json:"updatedat" sql:"updatedat"`
}

// Query operations

// Gets a specific data by id.
func (dt *Data) GetData(db *sql.DB) error {
	return db.QueryRow("SELECT strattr, intattr, createdat, updatedat FROM data WHERE id=$1",
		dt.ID).Scan(&dt.StrAttr, &dt.IntAttr, &dt.CreatedAt, &dt.UpdatedAt)
}

// Gets multiple data. Limit count and start position in db.
func GetMultipleData(db *sql.DB, start, count int) ([]Data, error) {
	rows, err := db.Query(
		"SELECT id, strattr, intattr, createdat, updatedat FROM data LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	data := []Data{}

	// Store query results into data variable if no errors.
	for rows.Next() {
		var dt Data
		if err := rows.Scan(&dt.ID, &dt.StrAttr, &dt.IntAttr, &dt.CreatedAt, &dt.UpdatedAt); err != nil {
			return nil, err
		}
		data = append(data, dt)
	}

	return data, nil
}

// CRUD operations

// Create new data and insert to database.
func (dt *Data) CreateData(db *sql.DB) error {
	// Scan db after creation if data exists using new data id.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO data(strattr, intattr, createdat, updatedat) VALUES($1, $2, $3, $4) RETURNING id, strattr, intattr, createdat, updatedat", dt.StrAttr, dt.IntAttr, timestamp, timestamp).Scan(&dt.ID, &dt.StrAttr, &dt.IntAttr, &dt.CreatedAt, &dt.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific data details by id.
func (dt *Data) UpdateData(db *sql.DB) error {
	timestamp := time.Now()
	_, err :=
		db.Exec("UPDATE data SET strattr=$1, intattr=$2, updatedat=$3 WHERE id=$4 RETURNING id, strattr, intattr, createdat, updatedat", dt.StrAttr, dt.IntAttr, timestamp, dt.ID)

	return err
}

// Deletes a specific data by id.
func (dt *Data) DeleteData(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM data WHERE id=$1", dt.ID)

	return err
}
