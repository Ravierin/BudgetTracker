package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateManualID generates a unique ID for manual entries
// Format: manual_<timestamp>_<random>
// This ensures no collision with exchange-generated IDs
func GenerateManualID(entityType string) string {
	timestamp := time.Now().UnixNano()
	
	// Generate 4 random bytes (8 hex chars)
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to timestamp only if random fails
		return fmt.Sprintf("manual_%s_%d", entityType, timestamp)
	}
	
	return fmt.Sprintf("manual_%s_%d_%x", entityType, timestamp, randomBytes)
}
