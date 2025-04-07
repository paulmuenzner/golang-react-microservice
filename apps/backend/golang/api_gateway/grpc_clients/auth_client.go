package grpc_clients

import (
	"log"

	"github.com/paulmuenzner/config"
	"github.com/paulmuenzner/protos" // Generated proto files for AuthService
	"google.golang.org/grpc"
)

// Set up the gRPC client for the AuthService
func GetAuthServiceClient() (protos.AuthServiceClient, error) {
	conn, err := grpc.Dial(config.AuthServiceURL, grpc.WithInsecure()) // Assuming gRPC is unsecured (use secure if needed)
	if err != nil {
		log.Fatalf("Failed to connect to AuthService: %v", err)
		return nil, err
	}

	client := protos.NewAuthServiceClient(conn)
	return client, nil
}
