package logging

import (
	"io"

	"github.com/rs/zerolog"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger/config"
)

// Default
type Default struct {
	// App Name
	AppName string
	// Environment
	Environment string
	// Host Name
	Hostname string
	// Client IP
	ClientIP string
}

// Logging
type Logging struct {
	// Configuration
	config config.Config
	// Rabbitmq
	RabbitMq domain.IRabbitMQ
	// SQS
	SQS domain.ISQS
	// App Name
	AppName string
	// Writers
	Writers []io.Writer
	// Log
	log zerolog.Logger
	// Default
	Default Default
}

type Type string

type Log struct {
	Type        Type             `json:"type"`
	SubType     domain.SubType   `json:"subType"`
	App         string           `json:"app"`
	Environment string           `json:"environment"`
	Error       string           `json:"error"`
	HostName    string           `json:"hostName"`
	ClientIp    string           `json:"clientIp"`
	Backend     *domain.Backend  `json:"backend,omitempty"`
	Frontend    *domain.Frontend `json:"frontend,omitempty"`
}
