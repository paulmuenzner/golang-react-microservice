package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	AppName     string `json:"appName"`
	Version     string `json:"version"`
	ApiEndpoint string `json:"apiEndpoint"`
}

// LoadConfig loads configuration from the baseConfig.json file
func LoadConfig(filePath string, configStruct interface{}) error {
	// Read the JSON file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading the config file: %v", err)
	}

	// Unmarshal the JSON data into the provided struct
	err = json.Unmarshal(data, configStruct)
	if err != nil {
		return fmt.Errorf("error unmarshalling config: %v", err)
	}

	return nil
}
