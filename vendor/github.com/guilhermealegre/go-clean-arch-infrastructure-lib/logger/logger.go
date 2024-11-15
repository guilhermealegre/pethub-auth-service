package logger

import (
	"io"
	"net"
	"os"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger/logging"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	loggerConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger/config"
)

// Logger service
type Logger struct {
	// Name
	name string
	// Configuration
	config *loggerConfig.Config
	// Writers
	writers []io.Writer
	// App
	app domain.IApp
	// Log
	log domain.ILogging
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	// configFile config file
	configFile = "logger.yaml"
)

// New logger service
func New(app domain.IApp, cfg *loggerConfig.Config, writers ...io.Writer) *Logger {
	logger := &Logger{
		app:     app,
		name:    "Logger",
		writers: writers,
	}

	if cfg != nil {
		logger.config = cfg
	}

	return logger
}

// Log gets the logger
func (l *Logger) Log() domain.ILogging {
	if l.log == nil {
		hostName, er := os.Hostname()
		if er != nil {
			hostName = "localhost"
		}
		l.log = &logging.Logging{
			RabbitMq: l.app.Rabbitmq(),
			SQS:      l.app.SQS(),
			AppName:  l.app.Config().Name,
			Default: logging.Default{
				AppName:     l.app.Config().Name,
				Environment: l.app.Config().Env,
				Hostname:    hostName,
				ClientIP:    GetClientIp(),
			},
			Writers: l.writers,
		}
		l.log.Init(*l.config)
	}

	return l.log
}

// Name service name
func (l *Logger) Name() string {
	return l.name
}

// Start starts the logger service
func (l *Logger) Start() (err error) {
	if l.config == nil {
		l.config = &loggerConfig.Config{}
		l.config.AdditionalConfig = l.additionalConfigType
		if err = config.Load(l.ConfigFile(), l.config); err != nil {
			err = errors.ErrorLoadingConfigFile().Formats(l.ConfigFile(), err)
			message.ErrorMessage(l.Name(), err)
			return err
		}
	}

	l.started = true

	return nil
}

// Stop stops the logger service
func (l *Logger) Stop() error {
	if !l.started {
		return nil
	}
	l.started = false
	return nil
}

// Config gets the configurations
func (l *Logger) Config() *loggerConfig.Config {
	return l.config
}

// ConfigFile gets the configuration file
func (l *Logger) ConfigFile() string {
	return configFile
}

// WithAdditionalConfigType sets an additional config type
func (l *Logger) WithAdditionalConfigType(obj interface{}) domain.ILogger {
	l.additionalConfigType = obj
	return l
}

// GetClientIp gets the client ip
func GetClientIp() (ip string) {
	addresses, _ := net.InterfaceAddrs()
	for _, a := range addresses {
		if ipNet, ok := a.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}
	return ip
}

func (l *Logger) DBLog(err error) error {
	if err == nil {
		return nil
	}
	l.Log().Do(err, &domain.LoggerInfo{
		SubType: domain.DatabaseSubType,
	})
	return errors.ErrorDatabase()
}

func (l *Logger) SQSLog(err error) error {
	if err == nil {
		return nil
	}
	l.Log().Do(err, &domain.LoggerInfo{
		SubType: domain.SQSSubType,
	})
	return errors.ErrorSQS()
}

func (l *Logger) ElasticLog(err error) error {
	if err == nil {
		return nil
	}
	l.Log().Do(err, &domain.LoggerInfo{
		SubType: domain.ElasticSubType,
	})
	return errors.ErrorElastic()
}

func (l *Logger) RedisLog(err error) error {
	if err == nil {
		return nil
	}
	l.Log().Do(err, &domain.LoggerInfo{
		SubType: domain.RedisSubType,
	})
	return errors.ErrorRedis()
}

// Started true if started
func (l *Logger) Started() bool {
	return l.started
}
