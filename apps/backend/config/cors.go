package config

import (
	"time"

	"github.com/gin-contrib/cors"
)

// CorsConfig returns a configured CORS middleware
func CorsConfig() cors.Config {
	return cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-frontend.com"}, // Allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},                       // Allowed HTTP methods
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},            // Allowed headers
		AllowCredentials: true,                                                           // Allow cookies and credentials
		MaxAge:           12 * time.Hour,                                                 // Cache the preflight response
	}
}
