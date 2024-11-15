package config

// Configs grpc configurations
type Configs struct {
	// Server
	Server Config `json:"server"`
	// Clients
	Clients []Config `json:"clients"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}

// Config configuration for server/client
type Config struct {
	// Name
	Name string `yaml:"name"`
	// Host
	Host string `yaml:"host"`
	// Port
	Port int `yaml:"port"`
}
