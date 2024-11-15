package logger

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/logger/config"
	"github.com/stretchr/testify/mock"
)

func NewLoggerMock() *LoggerMock {
	return &LoggerMock{}
}

type LoggerMock struct {
	mock.Mock
}

func (l *LoggerMock) Name() string {
	args := l.Called()
	return args.Get(0).(string)
}

func (l *LoggerMock) Start() error {
	args := l.Called()
	return args.Error(0)
}

func (l *LoggerMock) Stop() error {
	args := l.Called()
	return args.Error(0)
}

func (l *LoggerMock) ConfigFile() string {
	args := l.Called()
	return args.Get(0).(string)
}

func (l *LoggerMock) Config() *config.Config {
	args := l.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*config.Config)
}

func (l *LoggerMock) Log() domain.ILogging {
	args := l.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.ILogging)
}

func (l *LoggerMock) DBLog(err error) error {
	args := l.Called(err)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (l *LoggerMock) SQSLog(err error) error {
	args := l.Called(err)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (l *LoggerMock) ElasticLog(err error) error {
	args := l.Called(err)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func (l *LoggerMock) RedisLog(err error) error {
	args := l.Called(err)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

// WithAdditionalConfigType sets an additional config type
func (l *LoggerMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := l.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (l *LoggerMock) Started() bool {
	args := l.Called()
	return args.Get(0).(bool)
}
