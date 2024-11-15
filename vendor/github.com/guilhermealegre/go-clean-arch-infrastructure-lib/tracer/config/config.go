package config

import "strings"

// Config configurations for the tracer
type Config struct {
	// Enabled
	Enabled bool `yaml:"enabled"`
	// Log
	Log bool `yaml:"log"`
	// CollectorHostPort
	CollectorHostPort string `yaml:"collectorHostPort"`
	// Sensitive Uris
	SensitiveUris SensitiveUriList `yaml:"sensitiveUris"`
	// Additional Config
	AdditionalConfig interface{} `yaml:"additionalConfig"`
}

type SensitiveUriList []SensitiveUri
type SensitiveUri struct {
	Method string `yaml:"method"`
	Uri    string `yaml:"uri"`
}

func (list SensitiveUriList) Contains(method, uri string) bool {
	for _, i := range list {
		if strings.EqualFold(i.Method, method) &&
			strings.EqualFold(i.Uri, uri) {
			return true
		}
	}
	return false
}
