/*
   Basic CRUD functions for all the database operations
*/

package util

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Paste struct {
	ID        string
	Content   string
	CreatedAt time.Time
}

type DB struct {
	*sql.DB
}

func InitDB() (*DB, error) {

	db, err := NewDB("pastebin.db")
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if err := db.InitSchema(); err != nil {
		log.Fatalf("Failed to initialize schema: %v", err)
	}

	return db, err
}

func NewDB(dataSourceName string) (*DB, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &DB{db}, nil
}

func (db *DB) InitSchema() error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS pastes (
			id TEXT NOT NULL PRIMARY KEY,
			content TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema: %v", err)
	}
	return nil
}

func (db *DB) Create(paste *Paste) error {
	_, err := db.Exec("INSERT INTO pastes (id, content) VALUES (?, ?)", paste.ID, paste.Content)
	if err != nil {
		return fmt.Errorf("failed to create paste: %v", err)
	}
	return nil
}

func (db *DB) GetOne(id string) (*Paste, error) {
	paste := &Paste{}
	err := db.QueryRow("SELECT id, content, created_at FROM pastes WHERE id = ?", id).Scan(&paste.ID, &paste.Content, &paste.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("paste not found")
		}
		return nil, fmt.Errorf("failed to get paste: %v", err)
	}
	return paste, nil
}

func (db *DB) Update(paste *Paste) error {
	_, err := db.Exec("UPDATE pastes SET content = ? WHERE id = ?", paste.Content, paste.ID)
	if err != nil {
		return fmt.Errorf("failed to update paste: %v", err)
	}
	return nil
}

func (db *DB) Delete(id string) error {
	_, err := db.Exec("DELETE FROM pastes WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete paste: %v", err)
	}
	return nil
}
