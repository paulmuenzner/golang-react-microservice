// user_service.go

package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	config "github.com/paulmuenzner/api_gateway/config"
	"github.com/paulmuenzner/api_gateway/grpc_clients"
	"github.com/paulmuenzner/api_gateway/protos"
	ratelimit "github.com/paulmuenzner/shared/routes/ratelimit"
)

// SetupUserServiceRoutes sets up all routes for the User Service with its own middleware
func SetupUserServiceRoutes(r *gin.Engine) {
	// Load the security configuration
	securityConfig, err := config.LoadSecurityConfig()
	if err != nil {
		log.Fatalf("Error loading security config: %v", err)
	}

	// Convert the interval and cooldown duration to time.Duration
	interval := time.Duration(securityConfig.UserRateLimitInterval) * time.Second
	cooldownDuration := time.Duration(securityConfig.UserRateLimitCooldownDuration) * time.Second

	// Create a new rate limiter using values from the configuration
	rateLimiter := ratelimit.NewRateLimiter(securityConfig.UserRateLimitMaxRequests, interval, cooldownDuration)

	// Start a background goroutine to periodically reset IP request counters
	go rateLimiter.ResetRequests()

	// Create a route group for user-related endpoints
	userGroup := r.Group("/user")

	// Apply the rate limiter middleware only to the user group
	userGroup.Use(ratelimit.RateLimitMiddleware(rateLimiter))

	// Define the route(s)
	userGroup.GET("/:id", GetUserDetails)
}

// GetUserDetails handles GET /user/:id and forwards request to UserService via gRPC
func GetUserDetails(c *gin.Context) {
	// Create gRPC client for the user service
	client, err := grpc_clients.GetUserServiceClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to UserService"})
		return
	}

	// Extract user ID from path
	userID := c.Param("id")

	// Create gRPC request and fetch data
	req := &protos.GetUserRequest{UserId: userID}
	resp, err := client.GetUserDetails(c, req)
	if err != nil {
		log.Println("Error fetching user details:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}

	// Return the response as JSON
	c.JSON(http.StatusOK, resp)
}
