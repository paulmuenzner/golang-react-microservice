package security

import "github.com/google/uuid"

// GenerateID creates a new random UUID (v4) as string
// Returns a RFC 4122 compliant UUID suitable for request IDs, entity IDs, etc.
func GenerateID() string {
	return uuid.New().String()
}

// GenerateUUID creates a new random UUID (v4) and returns the UUID object
// Use this if you need the UUID type instead of just the string representation
func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// MustGenerateID creates a new UUID or panics on error
// Use only in initialization code where failure should stop the program
func MustGenerateID() string {
	return uuid.Must(uuid.NewRandom()).String()
}
