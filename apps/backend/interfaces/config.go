package interfaces

type BaseConfig struct {
	AppName         string                `json:"appName"`
	Version         string                `json:"version"`
	ApiEndpoint     string                `json:"apiEndpoint"`
	Cors            CorsConfig            `json:"cors"`
	SecurityHeaders SecurityHeadersConfig `json:"securityHeaders"`
}

// BaseConfig Nested objects
type CorsConfig struct {
	AllowOrigins []string `json:"allowOrigins"`
	AllowMethods []string `json:"allowMethods"`
	AllowHeaders []string `json:"allowHeaders"`
}

type SecurityHeadersConfig map[string]string
