package main

import (
	"backend/config"
	"backend/interfaces"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func setSecurityHeaders(c *gin.Context) {
	// Manually add security headers
	c.Header("Strict-Transport-Security", "max-age=31536000")                                       // HSTS for 1 year
	c.Header("X-Frame-Options", "DENY")                                                             // Prevent iframe embedding
	c.Header("X-XSS-Protection", "1; mode=block")                                                   // Enable XSS protection
	c.Header("X-Content-Type-Options", "nosniff")                                                   // Prevent MIME type sniffing
	c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'none'") // CSP
	c.Header("Referrer-Policy", "no-referrer")                                                      // Referrer policy
	c.Next()
}

// Globale MongoDB-Variable
var client *mongo.Client

func main() {
	////// CONFIG ///////////////////////////////
	baseConfig := interfaces.BaseConfig{}
	err := config.LoadConfig("../shared/data/baseConfig.json", &baseConfig)

	// Gin-Engine erstellen
	r := gin.Default()

	// CORS settings
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "https://your-frontend.com"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// Apply the custom middleware to set security headers
	r.Use(setSecurityHeaders)

	////// ENV //////////////////////////////////
	// Initialize configuration (load .env file)
	config.Env()
	// Use environment variables
	// appEnv := os.Getenv("APP_ENV")

	// Initialize the MongoDB connection
	config.InitMongoDB()

	// Sicherstellen, dass die Verbindung geschlossen wird, wenn die Anwendung endet
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal("Fehler beim Trennen der MongoDB-Verbindung:", err)
		}
		fmt.Println("Verbindung zur MongoDB erfolgreich getrennt.")
	}()

	// Apply the middleware globally to parse the request body
	r.Use(config.ParseRequestBodyMiddleware())

	// Einfacher Health-Check Endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Server starten auf Port 8080
	r.Run(":8080")
}
