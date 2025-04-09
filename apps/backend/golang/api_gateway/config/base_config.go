package config

import (
	"fmt"
	"path/filepath"

	file "github.com/paulmuenzner/shared/file"
)

// BaseConfig
type BaseConfig struct {
	AppName     string `json:"appName"`
	BackendPort string `json:"backendPort"`
	ApiEndpoint string `json:"apiEndpoint"`
}

func LoadBaseConfig() (*BaseConfig, error) {
	var baseCfg BaseConfig
	path := filepath.Join("..", "..", "..", "..", "..", "data/baseConfig.json")

	err := file.LoadJson(path, &baseCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %v", err)
	}
	return &baseCfg, nil
}
