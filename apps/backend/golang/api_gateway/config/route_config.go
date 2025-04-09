package config

import (
	"fmt"

	file "github.com/paulmuenzner/shared/file" // Importing the shared package
)

// RouteConfig holds gRPC service URLs
type RouteConfig struct {
	UserServiceURL string `json:"userServiceUrl"`
	AuthServiceURL string `json:"authServiceUrl"`
}

func LoadRouteConfig() (*RouteConfig, error) {
	var routeCfg RouteConfig
	err := file.LoadJson("data/route_config.json", &routeCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load route config: %v", err)
	}
	return &routeCfg, nil
}

// func main() {
// 	// Load the route config
// 	routeCfg, err := LoadRouteConfig()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Use the loaded config
// 	fmt.Printf("User Service URL: %s\n", routeCfg.UserServiceURL)
// 	fmt.Printf("Auth Service URL: %s\n", routeCfg.AuthServiceURL)
// }
