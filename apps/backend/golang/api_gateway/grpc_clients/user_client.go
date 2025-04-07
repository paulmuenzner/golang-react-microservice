package grpc_clients

import (
	"log"

	"github.com/yourrepo/config"
	"github.com/yourrepo/protos" // Generated proto files for UserService
	"google.golang.org/grpc"
)

// Set up the gRPC client for the UserService
func GetUserServiceClient() (protos.UserServiceClient, error) {
	conn, err := grpc.Dial(config.UserServiceURL, grpc.WithInsecure()) // Assuming gRPC is unsecured (use secure if needed)
	if err != nil {
		log.Fatalf("Failed to connect to UserService: %v", err)
		return nil, err
	}

	client := protos.NewUserServiceClient(conn)
	return client, nil
}
