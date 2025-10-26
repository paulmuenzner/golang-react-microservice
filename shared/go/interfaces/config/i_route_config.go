package interfaceconfig

// Sub-struct für einzelne Service-Routen
type IServiceRoute struct {
	Prefix    string `json:"Prefix"`
	TargetURL string `json:"TargetURL"`
}

// Struct for the new backend-specific fields
type IBackendConfig struct {
	ListenHost string        `json:"ListenHost"`
	ListenPort string        `json:"ListenPort"`
	ServiceA   IServiceRoute `json:"ServiceA"`
	ServiceB   IServiceRoute `json:"ServiceB"`
}

// Main config structure (must represent the "backend" level)
type IRouteConfig struct {
	Backend IBackendConfig `json:"backend"`
}
