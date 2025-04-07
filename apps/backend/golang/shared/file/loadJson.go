package file

import (
	"encoding/json"
	"os"
	"path/filepath"
)

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
