package config

import (
	"backend/interfaces"
	"log"
	"time"

	"github.com/gin-contrib/cors"
)

// CorsConfig returns a configured CORS middleware
func CorsConfig() cors.Config {
	var baseConfig interfaces.BaseConfig

	// Load the config
	err := LoadConfig("../shared/data/baseConfig.json", &baseConfig)
	if err != nil {
		log.Fatal("Failed to load baseConfig.json:", err)
	}

	// Extract CORS settings
	corsConfig := baseConfig.Cors

	return cors.Config{
		AllowOrigins:     corsConfig.AllowOrigins, // Allowed origins
		AllowMethods:     corsConfig.AllowMethods, // Allowed HTTP methods
		AllowHeaders:     corsConfig.AllowHeaders, // Allowed headers
		AllowCredentials: true,                    // Allow cookies and credentials
		MaxAge:           12 * time.Hour,          // Cache the preflight response
	}
}
