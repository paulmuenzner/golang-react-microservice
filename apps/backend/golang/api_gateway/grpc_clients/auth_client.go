package grpc_clients

import (
	"log"

	"github.com/paulmuenzner/api_gateway/config"
	"github.com/paulmuenzner/api_gateway/protos" // Generated proto files for AuthService
	config_shared "github.com/paulmuenzner/shared/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Set up the gRPC client for the AuthService
func GetAuthServiceClient() (protos.AuthServiceClient, error) {

	// Fetch environment variable and load configuration accordingly
	environment := config_shared.GetEnv("ENVIRONMENT", "local")

	// Get the AuthService URL based on the environment
	authServiceURL := config.GetServiceURL(environment, "AUTH_SERVICE_URL")

	// Konfiguration des Clients mit Transport-Anmeldeinformationen
	conn, err := grpc.NewClient(
		authServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Für unsichere Verbindungen
	)
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
		return nil, err
	}

	// Erstellung des Clients für den AuthService
	client := protos.NewAuthServiceClient(conn)
	return client, nil
}
