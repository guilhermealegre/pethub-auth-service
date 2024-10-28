package config

func NewConfig() *Config {
	return &Config{}
}

// Config redis configuration
type Config struct {
	// Host
	Host string `yaml:"host"`
	// Port
	Port int `yaml:"port"`
	// Password
	Password string `yaml:"password"`
	// Database
	Database int `yaml:"database"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}
