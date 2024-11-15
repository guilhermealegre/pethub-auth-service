package rabbitmq

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	rabbitmqConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/rabbitmq/config"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
)

func NewRabbitmqMock() *RabbitmqMock {
	return &RabbitmqMock{}
}

type RabbitmqMock struct {
	mock.Mock
}

func (r *RabbitmqMock) Name() string {
	args := r.Called()
	return args.Get(0).(string)
}

func (r *RabbitmqMock) Start() error {
	args := r.Called()
	return args.Error(0)
}

func (r *RabbitmqMock) Stop() error {
	args := r.Called()
	return args.Error(0)
}

func (r *RabbitmqMock) ConfigFile() string {
	args := r.Called()
	return args.Get(0).(string)
}

func (r *RabbitmqMock) Config() *rabbitmqConfig.Config {
	args := r.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*rabbitmqConfig.Config)
}

func (r *RabbitmqMock) Produce(message string, exchange string, routingKey string) error {
	args := r.Called(message, exchange, routingKey)
	return args.Error(0)
}

func (r *RabbitmqMock) Consume(app domain.IApp, queues string, handlers map[string]func(msg amqp.Delivery) bool) {
	_ = r.Called(app, queues, handlers)
}

func (r *RabbitmqMock) WithConsumer(consumer domain.IRabbitMQConsumer) domain.IRabbitMQ {
	args := r.Called(consumer)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(domain.IRabbitMQ)
}

// WithAdditionalConfigType sets an additional config type
func (r *RabbitmqMock) WithAdditionalConfigType(obj interface{}) domain.IApp {
	args := r.Called(obj)
	return args.Get(0).(domain.IApp)
}

// Started true if started
func (r *RabbitmqMock) Started() bool {
	args := r.Called()
	return args.Get(0).(bool)
}
