package config

import (
	"log"

	"github.com/paulmuenzner/shared/interfaces"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersMiddleware adds security-related HTTP headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	var baseConfig interfaces.BaseConfig

	// Load the config
	err := LoadConfig("../shared/data/baseConfig.json", &baseConfig)
	if err != nil {
		log.Fatal("Failed to load baseConfig.json:", err)
	}

	// Get security headers from config
	securityHeaders := baseConfig.SecurityHeaders

	return func(c *gin.Context) {
		for key, value := range securityHeaders {
			c.Header(key, value)
		}
		c.Next()
	}
}
