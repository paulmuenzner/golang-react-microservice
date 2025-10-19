package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	// logger "global/shared/logger"
)

// LoadConfig loads configuration from the baseConfig.json file
func LoadJson(filePath string, configStruct interface{}) error {
	// Read the JSON file
	// Build absolute path if needed
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		// logger.Error(err.Error())
		return err
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("error reading the file path: %v", err)
	}

	// Unmarshal the JSON data into the provided struct
	err = json.Unmarshal(data, configStruct)
	if err != nil {
		return fmt.Errorf("error unmarshalling json file: %v", err)
	}

	return nil
}
