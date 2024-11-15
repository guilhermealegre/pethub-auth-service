package database

import (
	"fmt"
	"github.com/gocraft/dbr/v2"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/database/middlewares/tracer"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"

	"github.com/guilhermealegre/go-clean-arch-core-lib/database/session"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	databaseConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/database/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

// Database has the database information
type Database struct {
	// Name
	name string
	// Configuration
	config *databaseConfig.Config
	// App
	app domain.IApp
	// Database
	client *database
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

// database connections
type database struct {
	// DbRead read connection
	DbRead session.ISession
	// DbWrite write connection
	DbWrite session.ISession
}

const (
	// configFile database config file name
	configFile = "database.yaml"
)

// New creates a new database
func New(app domain.IApp, config *databaseConfig.Config) *Database {
	database := &Database{
		name:   "Database",
		app:    app,
		client: &database{},
	}

	if config != nil {
		database.config = config
	}

	return database
}

// Name gets the database service name
func (d *Database) Name() string {
	return d.name
}

// Start starts the database connection
func (d *Database) Start() (err error) {
	if d.config == nil {
		d.config = &databaseConfig.Config{}
		d.config.AdditionalConfig = d.additionalConfigType
		if err = config.Load(d.ConfigFile(), d.config); err != nil {
			err = errors.ErrorLoadingConfigFile().Formats(d.ConfigFile(), err)
			message.ErrorMessage(d.Name(), err)
			return err
		}
	}

	// tracer
	middleware := tracer.NewTracerMiddleware(d.app)

	// master connection
	if d.client.DbWrite, err = openConnection(d.config.Master, middleware); err != nil {
		return err
	}

	// slave connection
	if d.client.DbRead, err = openConnection(d.config.Slave, middleware); err != nil {
		return err
	}

	migration := NewMigration(d.name, d.Write().DB(), d.config.Master.Driver, d.config.MigrationsDisabled)
	if errs := migration.Run(); len(errs) > 0 {
		for _, err = range errs {
			message.ErrorMessage(d.name, err)
		}
		return errors.ErrorMigration()
	}

	d.started = true

	return nil
}

// Config database configurations
func (d *Database) Config() *databaseConfig.Config {
	return d.config
}

// ConfigFile database configuration file
func (d *Database) ConfigFile() string {
	return configFile
}

// Read read connection
func (d *Database) Read() session.ISession {
	return d.client.DbRead
}

// Write write connection
func (d *Database) Write() session.ISession {
	return d.client.DbWrite
}

// openConnection open a new connection
func openConnection(config *databaseConfig.Connection, eventReceiver dbr.EventReceiver) (*session.Session, error) {
	conn, err := dbr.Open(
		config.Driver,
		fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s sslmode=disable port=%d",
			config.Host,
			config.User,
			config.Password,
			config.Database,
			config.Port,
		),
		eventReceiver,
	)

	if err != nil {
		return nil, err
	}

	return session.NewSession(conn, eventReceiver), nil
}

// Stop stops the database connection
func (d *Database) Stop() (err error) {

	if d.client.DbRead != nil {
		if err = d.client.DbRead.Close(); err != nil {
			return err
		}
	}

	if d.client.DbWrite != nil {
		if err = d.client.DbWrite.Close(); err != nil {
			return err
		}
	}

	d.started = false

	return nil
}

// WithAdditionalConfigType sets an additional config type
func (d *Database) WithAdditionalConfigType(obj interface{}) domain.IDatabase {
	d.additionalConfigType = obj
	return d
}

// Started true if started
func (d *Database) Started() bool {
	return d.started
}
