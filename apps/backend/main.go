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



// Globale MongoDB-Variable
var client *mongo.Client

func main() {
	////// CONFIG ///////////////////////////////
	baseConfig := interfaces.BaseConfig{}
	err := config.LoadConfig("../shared/data/baseConfig.json", &baseConfig)

	// Gin-Engine erstellen
	r := gin.Default()

	// 1️⃣ Apply CORS middleware (Handles cross-origin requests)
	r.Use(cors.New(config.CorsConfig()))

	// 2️⃣ Apply security headers middleware
	r.Use(config.SecurityHeadersMiddleware()) 

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
