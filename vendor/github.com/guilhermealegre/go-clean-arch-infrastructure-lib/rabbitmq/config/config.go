package config

// Config rabbitmq configurations
type Config struct {
	// User
	User string `yaml:"user"`
	// Password
	Password string `yaml:"password"`
	// Host
	Host string `yaml:"host"`
	// Port
	Port string `yaml:"port"`
	// VHost
	Vhost string `yaml:"vhost"`
	// Api
	Api string `yaml:"api"`
	// Api Header Host
	ApiHeaderHost string `yaml:"apiHeaderHost"`
	// Api Definitions
	ApiDefinitions string `yaml:"apiDefinitions"`
	// Queue Config File
	QueueConfigFile string `yaml:"queueConfigFile"`
	// Prefetch Count
	PrefetchCount int `yaml:"prefetchCount"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}
