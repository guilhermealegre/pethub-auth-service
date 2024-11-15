package config

// Config elastic search configurations
type Config struct {
	// Host
	Host string `yaml:"host"`
	// Port
	Port int `yaml:"port"`
	// User
	User string `yaml:"user"`
	// Password
	Password string `yaml:"password"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}
