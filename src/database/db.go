package database

import (
	"crypto/sha256"
	"database/sql"
	"ediscovery/src/models"
	"encoding/hex"
	"errors"
	"fmt"
)

import (
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Database represents the database connection
type Database struct {
	Conn *sql.DB
}

var database *Database

// NewDB creates a new database connection
func NewDB(dataSourceName string) error {
	conn, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return fmt.Errorf("failed to open the db: %w", err)
	}
	database = &Database{Conn: conn}

	err = createLegalholdTable()
	if err != nil {
		return fmt.Errorf("failed to createLegalholdTable: %w", err)
	}

	err = CreateUsersTable()
	if err != nil {
		return fmt.Errorf("failed to CreateUsersTable: %w", err)
	}

	return nil
}

// CreateUsersTable Function to create the users table
func CreateUsersTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        password TEXT NOT NULL
    );`

	_, err := database.Conn.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	return nil
}

func createLegalholdTable() error {
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

	_, err := database.Conn.Exec(createTableQuery)
	return err
}

// InsertKibana inserts a Kibana object into the database
func InsertKibana(k models.Kibana) error {
	insertQuery := `
	INSERT INTO legalhold (lhid, messageId, chatId, sent, sender, text, senderId, mediaType, groupName, participants)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err := database.Conn.Exec(insertQuery, k.LHID, k.ID, k.ChatID, k.Sent, k.Sender, k.Text, k.From, k.Type, k.Group, k.Participants)
	if err != nil {
		return fmt.Errorf("failed to insert kibana record: %w", err)
	}

	return nil
}

// FetchPaginatedChat returns a paginated list of messages for the given lhid and chatId
func FetchPaginatedChat(lhid string, chatId string, page int, limit int) ([]models.ChatMessage, error) {
	offset := (page - 1) * limit

	query := `
		SELECT sender, text, sent FROM legalhold
		WHERE lhid = ? AND chatId = ?
		ORDER BY sent DESC
		LIMIT ? OFFSET ?
	`
	rows, err := database.Conn.Query(query, lhid, chatId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var messages []models.ChatMessage
	for rows.Next() {
		var message models.ChatMessage
		if err := rows.Scan(&message.SenderName, &message.Text, &message.SentDate); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// FetchAllChats Fetch all chats for a given lhid
func FetchAllChats(lhid string) ([]models.Chat, error) {
	query := `SELECT chatId, groupName, participants FROM legalhold WHERE lhid = ? GROUP BY chatId`
	rows, err := database.Conn.Query(query, lhid)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	var chats []models.Chat
	for rows.Next() {
		var chat models.Chat
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

// FetchTotalMessagesCount returns the total number of messages for the given lhid and chatId
func FetchTotalMessagesCount(lhid string, chatId string) (int, error) {
	query := `
		SELECT COUNT(*) FROM legalhold
		WHERE lhid = ? AND chatId = ?
	`
	var count int
	err := database.Conn.QueryRow(query, lhid, chatId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// HashPassword Function to hash the password using SHA-256
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

// InsertUser Function to insert a new user (with hashed password)
func InsertUser(username, password string) error {
	hashedPassword := HashPassword(password)
	query := `INSERT INTO users (username, password) VALUES (?, ?);`
	_, err := database.Conn.Exec(query, username, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}
	return nil
}

// AuthenticateUser Function to authenticate user by username and password
func AuthenticateUser(username, password string) (bool, error) {
	var storedHash string
	query := `SELECT password FROM users WHERE username = ?;`
	err := database.Conn.QueryRow(query, username).Scan(&storedHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("invalid username or password")
		}
		return false, err
	}

	// Hash the input password and compare it with the stored hash
	inputHash := HashPassword(password)
	if inputHash != storedHash {
		return false, errors.New("invalid username or password")
	}
	return true, nil
}
