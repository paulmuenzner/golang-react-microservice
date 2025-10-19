package interfaceconfig

// Sub-struct für einzelne Service-Routen
type ServiceRoute struct {
	Prefix    string `json:"Prefix"`
	TargetURL string `json:"TargetURL"`
}

// Struct for the new backend-specific fields
type BackendConfig struct {
	ListenHost string       `json:"ListenHost"`
	ListenPort string       `json:"ListenPort"`
	ServiceA   ServiceRoute `json:"ServiceA"`
	ServiceB   ServiceRoute `json:"ServiceB"`
}

// Main config structure (must represent the "backend" level)
type RouteConfig struct {
	Backend BackendConfig `json:"backend"`
}
