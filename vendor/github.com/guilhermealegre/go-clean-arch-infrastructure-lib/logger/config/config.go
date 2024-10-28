package config

import "strings"

// Config logger configurations
type Config struct {
	// Path
	Path string `yaml:"path"`
	// File
	File string `yaml:"file"`
	// Level
	Level string `yaml:"level"`
	// StackTrace
	StackTrace bool `yaml:"stackTrace"`
	// Output
	Output []string `yaml:"output"`
	// Rabbitmq Config
	Rabbitmq *RabbitmqConfig `yaml:"rabbitmq"`
	// SQS Config
	SQS *SQSConfig `yaml:"sqs"`
	// Body
	Body bool `yaml:"body"`
	// Body Exclude Uris
	BodyExcludeUris BodyExcludeUriList `yaml:"bodyExcludeUris"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}

type BodyExcludeUriList []BodyExcludeUri
type BodyExcludeUri struct {
	Method string `yaml:"method"`
	Uri    string `yaml:"uri"`
}

type RabbitmqConfig struct {
	// Queue
	Queue string `yaml:"queue"`
	// Routing Key
	RoutingKey string `yaml:"routingKey"`
}

type SQSConfig struct {
	// Connection
	Connection string `yaml:"connection"`
	// Queue
	Queue string `yaml:"queue"`
	// Routing Key
	RoutingKey string `yaml:"routingKey"`
}

type MessageAttribute struct {
	DataType string `yaml:"dataType"`
	Value    string `yaml:"value"`
}

func (list BodyExcludeUriList) Contains(method, uri string) bool {
	for _, i := range list {
		if strings.EqualFold(i.Method, method) &&
			strings.EqualFold(i.Uri, uri) {
			return true
		}
	}
	return false
}
