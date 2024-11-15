package database

import (
	"github.com/guilhermealegre/go-clean-arch-core-lib/database/session"
	databaseConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/database/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/stretchr/testify/mock"
)

func NewDatabaseMock() *DatabaseMock {
	return &DatabaseMock{}
}

type DatabaseMock struct {
	mock.Mock
}

func (d *DatabaseMock) Name() string {
	args := d.Called()
	return args.Get(0).(string)
}

func (d *DatabaseMock) Start() error {
	args := d.Called()
	return args.Error(0)
}

func (d *DatabaseMock) Stop() error {
	args := d.Called()
	return args.Error(0)
}

func (d *DatabaseMock) ConfigFile() string {
	args := d.Called()
	return args.Get(0).(string)
}

func (d *DatabaseMock) Config() *databaseConfig.Config {
	args := d.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*databaseConfig.Config)
}

func (d *DatabaseMock) Read() session.ISession {
	args := d.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(session.ISession)
}

func (d *DatabaseMock) Write() session.ISession {
	args := d.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(session.ISession)
}

// WithAdditionalConfigType sets an additional config type
func (d *DatabaseMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := d.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (d *DatabaseMock) Started() bool {
	args := d.Called()
	return args.Get(0).(bool)
}
