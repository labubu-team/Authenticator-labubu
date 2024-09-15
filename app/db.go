package app

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // Import SQLite3 driver
	"os"
	"path/filepath"
)

// Database struct
type Database struct {
	db *sql.DB
}

// NewDatabase
func NewDatabase(dbName string) (*Database, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %w", err)
	}

	appSupportDir := filepath.Join(homeDir, "Library", "Application Support", "2FA-PHP")
	if err := os.MkdirAll(appSupportDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("error creating application support directory: %w", err)
	}

	dbPath := filepath.Join(appSupportDir, dbName)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS two_fa (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		priority INTEGER,
		logo TEXT,
		name TEXT,
		secret TEXT,
		domain TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}
