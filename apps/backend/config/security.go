package config

import "github.com/gin-gonic/gin"

// SecurityHeadersMiddleware adds security-related HTTP headers to responses
func SecurityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Strict-Transport-Security", "max-age=31536000")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none'")
		c.Header("Referrer-Policy", "no-referrer")
		c.Next()
	}
}
