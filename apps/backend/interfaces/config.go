// /shared/config/config.go
package interfaces

type BaseConfig struct {
	AppName     string `json:"appName"`
	Version     string `json:"version"`
	ApiEndpoint string `json:"apiEndpoint"`
}
