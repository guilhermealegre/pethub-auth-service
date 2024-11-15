package config

// Config S3 configurations
type Config struct {
	// Region
	Region string `yaml:"region"`
	// Bucket
	Bucket string `yaml:"bucket"`
	// Role
	Role string `yaml:"role"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}
