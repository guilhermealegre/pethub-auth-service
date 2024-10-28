package config

// Config configurations for the app
type Config struct {
	// Environment
	Env string `yaml:"env"`
	// Name
	Name string `yaml:"name"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}
