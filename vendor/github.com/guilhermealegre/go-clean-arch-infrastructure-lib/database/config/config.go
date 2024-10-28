package config

// Config configurations for the database
type Config struct {
	// Multi Database
	MultiDb bool `yaml:"multiDb"`
	// Master
	Master *Connection `yaml:"master"`
	// Slave
	Slave *Connection `yaml:"slave"`
	// Migrations Disabled
	MigrationsDisabled bool `yaml:"migrationsDisabled"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}

// Connection has the connection configuration
type Connection struct {
	// Database
	Database string `yaml:"database"`
	// Host
	Host string `yaml:"host"`
	// User
	User string `yaml:"user"`
	// Password
	Password string `yaml:"password"`
	// Port
	Port int `yaml:"port"`
	// Driver
	Driver string `yaml:"driver"`
}
