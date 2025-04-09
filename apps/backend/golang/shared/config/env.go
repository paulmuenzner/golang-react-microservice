package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Env loads environment variables from the .env file globally
func LoadEnv() error {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return nil
}

// GetEnv retrieves the value of an environment variable.
// If the variable is not set, it returns the default value.
func GetEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		// If the environment variable is not set, return the default value
		return defaultValue
	}
	return value
}
