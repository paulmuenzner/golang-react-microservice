package config

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client // Global MongoDB client variable

// InitMongoDB establishes a connection to the MongoDB database.
func InitMongoDB() {
	// MongoDB-Verbindungs-URI (pass with your MongoDB URI)
	uri := "mongodb://localhost:27017" // Example URI, change as needed

	// Client options
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB
	var err error
	Client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping MongoDB to ensure the connection works
	err = Client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	} else {
		fmt.Println("Successfully connected to MongoDB!")
	}
}
