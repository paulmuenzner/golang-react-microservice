package main

import (
	"backend/config"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// Globale MongoDB-Variable
var client *mongo.Client

func main() {
	// Initialize the MongoDB connection
	config.InitMongoDB()

	// Sicherstellen, dass die Verbindung geschlossen wird, wenn die Anwendung endet
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal("Fehler beim Trennen der MongoDB-Verbindung:", err)
		}
		fmt.Println("Verbindung zur MongoDB erfolgreich getrennt.")
	}()

	// Gin-Engine erstellen
	r := gin.Default()

	// Einfacher Health-Check Endpoint
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Server starten auf Port 8080
	r.Run(":8080")
}
