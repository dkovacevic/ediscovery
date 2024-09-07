package main

import (
	"database/sql"
	"fmt"
)

import (
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// DB represents the database connection
type DB struct {
	Conn *sql.DB
}

type ChatMessage struct {
	SenderName string
	Text       string
	SentDate   string
}

// Chat represents a single chat overview
type Chat struct {
	ChatID       string `json:"ChatID"`
	GroupName    string `json:"GroupName"`
	Participants string `json:"Participants"`
}

// NewDB creates a new database connection
func NewDB(dataSourceName string) (*DB, error) {
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Create the Kibana table if it doesn't exist
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS legalhold (
		lhid TEXT,
		messageId TEXT,
		chatId TEXT,
		sent TEXT,
		sender TEXT,
		text TEXT,
		senderId TEXT,
		mediaType TEXT,
		groupName TEXT,
		participants TEXT
	);`

	_, err = conn.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	return &DB{Conn: conn}, nil
}

// InsertKibana inserts a Kibana object into the database
func (db *DB) InsertKibana(k Kibana) error {
	insertQuery := `
	INSERT INTO legalhold (lhid, messageId, chatId, sent, sender, text, senderId, mediaType, groupName, participants)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := db.Conn.Exec(insertQuery, k.LHID, k.ID, k.ChatID, k.Sent, k.Sender, k.Text, k.From, k.Type, k.Group, k.Participants)
	if err != nil {
		return fmt.Errorf("failed to insert kibana record: %w", err)
	}

	return nil
}

// Fetch chat messages based on lhid and chatid
func (db *DB) fetchChat(lhid, chatid string) ([]ChatMessage, error) {
	query := `SELECT sender, text, sent FROM legalhold WHERE lhid = ? AND chatId = ? ORDER BY sent`
	rows, err := db.Conn.Query(query, lhid, chatid)
	if err != nil {
		return nil, err
	}

	var messages []ChatMessage
	for rows.Next() {
		var message ChatMessage
		if err := rows.Scan(&message.SenderName, &message.Text, &message.SentDate); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

// Fetch all chats for a given lhid
func (db *DB) fetchAllChats(lhid string) ([]Chat, error) {
	query := `SELECT chatId, groupName, participants FROM legalhold WHERE lhid = ? GROUP BY chatId`
	rows, err := db.Conn.Query(query, lhid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ChatID, &chat.GroupName, &chat.Participants); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chats, nil
}
