package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const createTableQuery = `
CREATE TABLE IF NOT EXISTS compromised_emails (
    id INTEGER PRIMARY KEY,
    email_hash VARCHAR(64) UNIQUE,
    breach_date TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

func InitDatabase(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully")
	return db, nil
}