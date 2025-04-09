package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/paulmuenzner/api_gateway/config"
	"github.com/paulmuenzner/api_gateway/routes"
	"github.com/paulmuenzner/shared/date"
	logger "github.com/paulmuenzner/shared/logging"
)

func main() {
	// Initialize logger
	logger.Info("Start api_gateway now")

	// Load the route base
	baseCfg, err := config.LoadBaseConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Setting up UserService routes")

	// Initialize Gin router
	r := gin.Default()

	// Setup routes for the API Gateway
	routes.SetupUserServiceRoutes(r)
	routes.SetupAuthServiceRoutes(r)
	routes.SetupLoggingServiceRoutes(r)

	// Start server
	port := baseCfg.BackendPort
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}

	// Print current time to ensure everything is running
	currentTime := date.GetCurrentUTCTime()
	fmt.Println("Current time:", date.FormatDate(currentTime))
}
