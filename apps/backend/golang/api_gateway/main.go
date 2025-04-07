package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paulmuenzner/api_gateway/routes"
	"github.com/paulmuenzner/shared/date"
	logger "github.com/paulmuenzner/shared/logging"
)

func main() {
	// Initialize logger
	logger.Info("Start api_gateway")

	// Initialize Gin router
	r := gin.Default()

	// Setup routes for the API Gateway
	routes.SetupUserServiceRoutes(r)
	routes.SetupAuthServiceRoutes(r)

	// Start server
	port := "8080" // You can use environment variables or a config file
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}

	// Print current time to ensure everything is running
	currentTime := date.GetCurrentUTCTime()
	fmt.Println("Current time:", date.FormatDate(currentTime))
}
