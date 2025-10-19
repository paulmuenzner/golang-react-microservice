package config

import (
	"bytes"
	"encoding/json"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

// Middleware to parse the request body globally
func ParseRequestBodyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Read the entire request body
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println("Error reading request body:", err)
			c.JSON(400, gin.H{"error": "Cannot read request body"})
			c.Abort()
			return
		}

		// Re-set the body to allow other handlers to read it if needed
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// Use a map to store any kind of JSON object
		var parsedBody map[string]interface{}
		err = json.Unmarshal(body, &parsedBody)
		if err != nil {
			log.Println("Error parsing request body:", err)
			c.JSON(400, gin.H{"error": "Invalid JSON"})
			c.Abort()
			return
		}

		// Store the parsed body in the context for later use in routes
		c.Set("parsedBody", parsedBody)

		// Proceed to the next middleware/handler
		c.Next()
	}
}
