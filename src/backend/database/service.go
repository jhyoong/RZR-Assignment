package database

import (
	"database/sql"
	"time"
)

type EmailService struct {
	db *sql.DB
}

func NewEmailService(db *sql.DB) *EmailService {
	return &EmailService{db: db}
}

func (s *EmailService) IsEmailCompromised(emailHash string) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM compromised_emails WHERE email_hash = ?"
	err := s.db.QueryRow(query, emailHash).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *EmailService) AddCompromisedEmail(emailHash string, breachDate time.Time) error {
	query := "INSERT OR IGNORE INTO compromised_emails (email_hash, breach_date) VALUES (?, ?)"
	_, err := s.db.Exec(query, emailHash, breachDate)
	return err
}

func (s *EmailService) GetCompromisedEmailCount() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM compromised_emails"
	err := s.db.QueryRow(query).Scan(&count)
	return count, err
}