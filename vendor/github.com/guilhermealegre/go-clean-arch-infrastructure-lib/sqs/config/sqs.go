package config

// Config sqs configurations
type Config struct {
	// Connections
	Connections map[string]*Connection `yaml:"connections"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}

type Connection struct {
	// Credentials
	Credentials Credentials `yaml:"credentials"`
	// Migrations Disabled
	MigrationsDisabled bool `yaml:"migrationsDisabled"`
	// Max Number of Messages
	MaxNumberOfMessages int64 `yaml:"maxNumberOfMessages"`
	// Prefetch Count
	VisibilityTimeout int64 `yaml:"visibilityTimeOut"`
	// Wait Time Seconds
	WaitTimeSeconds int64 `yaml:"waitTimeSeconds"`
	// Add Environment Prefix Queue
	AddEnvPrefixQueue bool `yaml:"addEnvPrefixQueue"`
}

type Credentials struct {
	// Region
	Region string `yaml:"region"`
	// Api
	Api string `yaml:"api"`
	// Id Account
	IdAccount string `yaml:"idAccount"`
}
