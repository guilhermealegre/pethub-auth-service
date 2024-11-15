package redis

import (
	"fmt"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/redis/middlewares/tracer"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
	redisConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/redis/config"
	"github.com/redis/go-redis/v9"
)

// Redis service
type Redis struct {
	// Name
	name string
	// Configuration
	config *redisConfig.Config
	// App
	app domain.IApp
	// Client
	client *redis.Client
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	// configFile redis configuration file
	configFile = "redis.yaml"
)

// New creates a new redis service
func New(app domain.IApp, config *redisConfig.Config) *Redis {
	redis := &Redis{
		name: "Redis",
		app:  app,
	}

	if config != nil {
		redis.config = config
	}

	return redis
}

// Name gets the service name
func (r *Redis) Name() string {
	return r.name
}

// Start starts the redis service
func (r *Redis) Start() (err error) {
	if r.config == nil {
		r.config = redisConfig.NewConfig()
		r.config.AdditionalConfig = r.additionalConfigType
		if err = config.Load(r.ConfigFile(), r.config); err != nil {
			err = errors.ErrorLoadingConfigFile().Formats(r.ConfigFile(), err)
			message.ErrorMessage(r.Name(), err)
			return err
		}
	}

	address := r.config.Host
	if r.config.Port > 0 {
		address += fmt.Sprintf(":%d", r.config.Port)
	}

	r.client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: r.config.Password,
		DB:       r.config.Database,
	})

	// tracer
	r.client.AddHook(tracer.NewTracerMiddleware(r.app))

	r.started = true

	return nil
}

// Stop stops the redis service
func (r *Redis) Stop() error {
	if !r.started {
		return nil
	}
	r.started = false
	return r.client.Close()
}

// Config gets the configurations
func (r *Redis) Config() *redisConfig.Config {
	return r.config
}

// Client gets the redis client
func (r *Redis) Client() *redis.Client {
	return r.client
}

// ConfigFile gets the configuration file
func (r *Redis) ConfigFile() string {
	return configFile
}

// WithAdditionalConfigType sets an additional config type
func (r *Redis) WithAdditionalConfigType(obj interface{}) domain.IRedis {
	r.additionalConfigType = obj
	return r
}

// Started true if started
func (r *Redis) Started() bool {
	return r.started
}
