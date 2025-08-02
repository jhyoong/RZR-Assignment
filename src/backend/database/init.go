package database

import (
	"database/sql"
	"log"
	"time"

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
	// Enable WAL mode and other optimizations for concurrent access
	dsn := dbPath + "?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=1000&_foreign_keys=true"
	
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Set connection pool settings for better concurrency
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if _, err := db.Exec(createTableQuery); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully with WAL mode")
	return db, nil
}