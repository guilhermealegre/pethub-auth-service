package redis

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	redisConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/redis/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
)

func NewRedisMock() *RedisMock {
	return &RedisMock{}
}

type RedisMock struct {
	mock.Mock
}

func (r *RedisMock) Name() string {
	args := r.Called()
	return args.Get(0).(string)
}

func (r *RedisMock) Start() error {
	args := r.Called()
	return args.Error(0)
}

func (r *RedisMock) Stop() error {
	args := r.Called()
	return args.Error(0)
}

func (r *RedisMock) ConfigFile() string {
	args := r.Called()
	return args.Get(0).(string)
}

func (r *RedisMock) Config() *redisConfig.Config {
	args := r.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redisConfig.Config)
}

func (r *RedisMock) Client() *redis.Client {
	args := r.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*redis.Client)
}

// WithAdditionalConfigType sets an additional config type
func (r *RedisMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := r.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (r *RedisMock) Started() bool {
	args := r.Called()
	return args.Get(0).(bool)
}
