package config

import (
	"fmt"
	"log"

	file "github.com/paulmuenzner/shared/file" // Importing the shared package
)

// RouteConfig holds gRPC service URLs
type RouteConfig struct {
	Local map[string]string `json:"local"`
	K8s   map[string]string `json:"k8s"`
}

func LoadRouteConfig() (*RouteConfig, error) {
	var routeCfg RouteConfig
	err := file.LoadJson("data/route_config.json", &routeCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load route config: %v", err)
	}
	return &routeCfg, nil
}

// GetServiceURL fetches the URL for a service based on the environment and service name
func GetServiceURL(environment string, serviceName string) string {
	// Default environment to 'local' if not provided
	if environment == "" {
		environment = "local"
	}

	// Load the route configuration
	routeCfg, err := LoadRouteConfig()
	if err != nil {
		log.Fatalf("Error loading route config: %v", err)
	}

	var configForEnv map[string]string
	if environment == "local" {
		configForEnv = routeCfg.Local
	} else if environment == "k8s" {
		configForEnv = routeCfg.K8s
	}

	// Retrieve the URL for the given service
	serviceURL, ok := configForEnv[serviceName]
	if !ok {
		log.Fatalf("Service %s not found in config", serviceName)
	}

	return serviceURL
}
