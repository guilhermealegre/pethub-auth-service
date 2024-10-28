package config

// Config http configurations
type Config struct {
	// Host
	Host string `yaml:"host"`
	// Port
	Port int `yaml:"port"`
	// JWT Secret
	JwtSecret string `yaml:"jwtSecret"`
	// Jwt Expiry Time Hours
	JwtExpiryTimeHours int `yaml:"jwtExpiryTimeHours"`
	// Cookie Inactivity Minutes
	CookieInactivityMinutes int `yaml:"cookieInactivityMinutes"`
	// Api Keys
	ApiKeys []string `yaml:"apiKeys"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}
