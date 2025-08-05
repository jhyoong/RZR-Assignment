package database

import (
	"crypto/sha256"
	"fmt"
	"time"
)

func SeedDatabase(service *EmailService) error {
	testEmails := []string{
		"test@example.com",
		"user@domain.com",
		"compromised@email.com",
		"breach@test.org",
		"victim@company.co",
	}

	for _, email := range testEmails {
		emailHash := hashEmail(email)
		breachDate := time.Now().AddDate(0, 0, -30) // 30 days ago

		err := service.AddCompromisedEmail(emailHash, breachDate)
		if err != nil {
			return fmt.Errorf("failed to add %s: %w", email, err)
		}
	}

	return nil
}

func hashEmail(email string) string {
	hash := sha256.Sum256([]byte(email))
	return fmt.Sprintf("%x", hash)
}
