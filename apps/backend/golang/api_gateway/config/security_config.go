// security_config.go

package config

import (
	"fmt"

	file "github.com/paulmuenzner/shared/file"
)

// SecurityConfig holds gRPC service URLs
type SecurityConfig struct {
	UserRateLimitMaxRequests      int `yaml:"userRateLimitMaxRequests"`      // Max requests
	UserRateLimitInterval         int `yaml:"userRateLimitInterval"`         // Interval in seconds
	UserRateLimitCooldownDuration int `yaml:"userRateLimitCooldownDuration"` // Cooldown in seconds
}

func LoadSecurityConfig() (*SecurityConfig, error) {
	var securityCfg SecurityConfig
	err := file.LoadYAML("data/security_config.json", &securityCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load security config: %v", err)
	}
	return &securityCfg, nil
}
