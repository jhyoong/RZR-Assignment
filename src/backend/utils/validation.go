package utils

import (
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	return emailRegex.MatchString(email)
}

func SanitizeEmail(email string) string {
	// Convert to lowercase and trim whitespace
	email = strings.ToLower(strings.TrimSpace(email))

	// Remove any null bytes or control characters
	email = strings.ReplaceAll(email, "\x00", "")

	return email
}

func HashEmail(email string) string {
	// First sanitize the email
	sanitized := SanitizeEmail(email)

	// Create SHA-256 hash
	hash := sha256.Sum256([]byte(sanitized))
	return fmt.Sprintf("%x", hash)
}

func ValidateAndHashEmail(email string) (string, bool) {
	sanitized := SanitizeEmail(email)

	if !IsValidEmail(sanitized) {
		return "", false
	}

	return HashEmail(sanitized), true
}
