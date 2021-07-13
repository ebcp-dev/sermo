package model

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

// Defines channel model.
type Channel struct {
	ChannelID     uuid.UUID `json:"channelid" sql:"uuid"`
	ChannelName   string    `json:"channelname" validate:"required"`
	MaxPopulation int       `json:"maxpopulation" validate:"required"`
	UserID        uuid.UUID `json:"userid" sql:"uuid"`
	CreatedAt     time.Time `json:"createdat" validate:"required"`
	UpdatedAt     time.Time `json:"updatedat" validate:"required"`
}

// Query operations

// Gets a specific channel by ChannelID.
func (ch *Channel) GetChannel(db *sql.DB) error {
	return db.QueryRow("SELECT channelname, maxpopulation, createdat, updatedat FROM channels WHERE channelid=$1",
		ch.ChannelID).Scan(&ch.ChannelName, &ch.MaxPopulation, &ch.CreatedAt, &ch.UpdatedAt)
}

// Gets multiple channel. Limit count and start position in db.
func GetChannels(db *sql.DB, start, count int) ([]Channel, error) {
	rows, err := db.Query(
		"SELECT channelid, channelname, maxpopulation, userid, createdat, updatedat FROM channels LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}
	// Wait for query to execute then close the row.
	defer rows.Close()

	channel := []Channel{}

	// Store query results into channel variable if no errors.
	for rows.Next() {
		var ch Channel
		if err := rows.Scan(&ch.ChannelID, &ch.ChannelName, &ch.MaxPopulation, &ch.UserID, &ch.CreatedAt, &ch.UpdatedAt); err != nil {
			return nil, err
		}
		channel = append(channel, ch)
	}

	return channel, nil
}

// CRUD operations

// Create new channel and insert to database.
func (ch *Channel) CreateChannel(db *sql.DB) error {
	// Scan db after creation if channel exists using new channel ChannelID.
	timestamp := time.Now()
	err := db.QueryRow(
		"INSERT INTO channels(channelname, maxpopulation, userid, createdat, updatedat) VALUES($1, $2, $3, $4, $5) RETURNING channelid, channelname, maxpopulation, userid, createdat, updatedat", ch.ChannelName, ch.MaxPopulation, ch.UserID, timestamp, timestamp).Scan(&ch.ChannelID, &ch.ChannelName, &ch.MaxPopulation, &ch.UserID, &ch.CreatedAt, &ch.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Updates a specific channel details by ChannelID.
func (ch *Channel) UpdateChannel(db *sql.DB) error {
	timestamp := time.Now()
	err :=
		db.QueryRow("UPDATE channels SET channelname=$1, maxpopulation=$2, updatedat=$3 WHERE channelid=$4 RETURNING channelid, channelname, maxpopulation, userid, createdat, updatedat", ch.ChannelName, ch.MaxPopulation, timestamp, ch.ChannelID).Scan(&ch.ChannelID, &ch.ChannelName, &ch.MaxPopulation, &ch.UserID, &ch.CreatedAt, &ch.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

// Deletes a specific channel by ChannelID.
func (ch *Channel) DeleteChannel(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM channels WHERE channelid=$1", ch.ChannelID)
	if err != nil {
		return err
	}

	return nil
}
