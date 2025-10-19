package config

import (
	"fmt"
	"path/filepath"
	"runtime"
	"sync"

	interfaceconfig "github.com/app/shared/go/interfaces/config"
	file "github.com/app/shared/go/utils/files"
)

// getCallerDir for receiving absolute path of the caller directory,
func getCallerDir() (string, error) {
	// 0 = getCallerDir, 1 = LoadRelativeConfig
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("could not get the caller file path")
	}
	return filepath.Dir(filename), nil
}

// loadRelativeConfig is the generic, one-time loading function.
// configPath: The path to the JSON file, relative to the package's root directory.
// configStruct: A pointer to the target structure.
func loadRelativeConfig(configPath string, configStruct interface{}) error {
	callerDir, err := getCallerDir()
	if err != nil {
		return err
	}

	// Resolves the relative path (e.g., "../../data/config/baseConfig.json")
	// based on the anchor point (callerDir).
	fullPath := filepath.Join(callerDir, configPath)

	// Cleans the path to resolve '..' and get the final absolute path
	absolutePath := filepath.Clean(fullPath)

	// Calls the actual loading logic
	return file.LoadJson(absolutePath, configStruct)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////////////////////////////
// The following functions are wrapper for loading the routing configuration.
// They ensure that the configuration is loaded only once and makes the result accessible.

// ///////////////////////////
// Route configurations
var (
	routeConfigOnce sync.Once
	RouteConfig     interfaceconfig.RouteConfig
	routeConfigErr  error
)

func LoadRouteConfig() (interfaceconfig.RouteConfig, error) {
	routeConfigOnce.Do(func() {
		const relativePath = "../../data/config/routeConfig.json"
		routeConfigErr = loadRelativeConfig(relativePath, &RouteConfig)
	})
	return RouteConfig, routeConfigErr
}
