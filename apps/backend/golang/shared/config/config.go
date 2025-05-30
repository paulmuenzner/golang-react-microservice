package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	AppName     string `json:"appName"`
	Version     string `json:"version"`
	ApiEndpoint string `json:"apiEndpoint"`
}

// Load json
func LoadJson(filePath string, configStruct interface{}) error {
	// Read the JSON file
	// Build absolute path if needed
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, configStruct); err != nil {
		return err
	}

	return nil
}
