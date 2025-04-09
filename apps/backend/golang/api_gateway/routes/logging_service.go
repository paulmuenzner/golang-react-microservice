// logging_service.go

package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paulmuenzner/api_gateway/grpc_clients"
	"github.com/paulmuenzner/api_gateway/protos" // Generated proto files for UserService
)

func SetupLoggingServiceRoutes(r *gin.Engine) {
	r.GET("/user/:id", GetUserDetails)
}

func PostLoggingDetails(c *gin.Context) {
	// Get the gRPC client for UserService
	client, err := grpc_clients.GetUserServiceClient()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to UserService"})
		return
	}

	// Make a gRPC request to fetch user details
	userID := c.Param("id")
	req := &protos.GetUserRequest{UserId: userID}
	resp, err := client.GetUserDetails(c, req)
	if err != nil {
		log.Println("Error fetching user details:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user details"})
		return
	}

	// Return the user details
	c.JSON(http.StatusOK, resp)
}
