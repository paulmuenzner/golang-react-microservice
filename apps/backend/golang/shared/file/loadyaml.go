package file

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// LoadYAML loads a YAML file into a given config structure
func LoadYAML(filePath string, configStruct interface{}) error {
	// Build absolute path if needed
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}

	// Read the YAML file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return err
	}

	// Unmarshal the YAML data into the config structure
	if err := yaml.Unmarshal(data, configStruct); err != nil {
		return err
	}

	return nil
}
