package identity

import "github.com/google/uuid"

// NewPublicID generates a public, API-safe identifier (UUID v7)
func NewPublicID() string {
	return uuid.Must(uuid.NewV7()).String()
}
