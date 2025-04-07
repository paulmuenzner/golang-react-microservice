package config

import (
	"log"

	"github.com/joho/godotenv"
)

// Env loads environment variables from the .env file globally
func Env() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
