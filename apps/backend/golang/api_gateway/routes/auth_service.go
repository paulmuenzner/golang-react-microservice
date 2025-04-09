// package routes

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/paulmuenzner/api_gateway/grpc_clients"
// 	"github.com/paulmuenzner/api_gateway/protos" // Generated proto files for AuthService
// )

// func SetupAuthServiceRoutes(r *gin.Engine) {
// 	r.POST("/auth/login", Login)
// }

// func Login(c *gin.Context) {
// 	// Get the gRPC client for AuthService
// 	client, err := grpc_clients.GetAuthServiceClient()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to AuthService"})
// 		return
// 	}

// 	// Make a gRPC request to perform login
// 	var loginReq protos.LoginRequest
// 	if err := c.ShouldBindJSON(&loginReq); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	resp, err := client.Login(c, &loginReq)
// 	if err != nil {
// 		log.Println("Login failed:", err)
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
// 		return
// 	}

// 	// Return the auth token
// 	c.JSON(http.StatusOK, gin.H{"auth_token": resp.GetAuthToken()})
// }
