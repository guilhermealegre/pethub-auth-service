package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/config"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
	rabbitmqConfig "github.com/guilhermealegre/go-clean-arch-infrastructure-lib/rabbitmq/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Rabbitmq service
type Rabbitmq struct {
	// Name
	name string
	// App
	app domain.IApp
	// Configuration
	config *rabbitmqConfig.Config
	// Consumer Connection
	consumerConnection *amqp.Connection
	// Producer Connection
	producerConnection *amqp.Connection
	// Consumer Channel
	consumerChannel *amqp.Channel
	// Producer Channel
	producerChannel *amqp.Channel
	// Consumers
	consumers []domain.IRabbitMQConsumer
	// Additional Config Type
	additionalConfigType interface{}
	// Started
	started bool
}

const (
	// configFile rabbitmq configuration file
	configFile = "rabbitmq.yaml"
)

// New creates a new rabbitmq service
func New(app domain.IApp, config *rabbitmqConfig.Config) *Rabbitmq {
	rabbitmq := &Rabbitmq{
		app:  app,
		name: "Rabbitmq",
	}

	if config != nil {
		rabbitmq.config = config
	}

	return rabbitmq
}

// Name gets the service name
func (r *Rabbitmq) Name() string {
	return r.name
}

// WithConsumer adds a new consumer
func (r *Rabbitmq) WithConsumer(consumer domain.IRabbitMQConsumer) domain.IRabbitMQ {
	r.consumers = append(r.consumers, consumer)
	return r
}

// Start starts the rabbitmq service
func (r *Rabbitmq) Start() (err error) {
	if r.config == nil {
		r.config = &rabbitmqConfig.Config{}
		r.config.AdditionalConfig = r.additionalConfigType
		// load the configuration file
		if err = config.Load(r.ConfigFile(), r.config); err != nil {
			err = errors.ErrorLoadingConfigFile().Formats(r.ConfigFile(), err)
			message.ErrorMessage(r.Name(), err)
			return err
		}
	}

	r.Connect()

	if err = r.runMigration(); err != nil {
		message.ErrorMessage(r.Name(), err)
		return err
	}

	//load the migration file
	for _, handler := range r.consumers {
		r.Consume(r.app, handler.GetQueue(), handler.GetHandlers())
	}

	r.started = true

	return nil
}

// Connect connects the rabbitmq
func (r *Rabbitmq) Connect() {
	var err error

	amqpInfo := fmt.Sprintf("amqp://%s:%s@%s:%s%s", r.config.User, r.config.Password, r.config.Host, r.config.Port, r.config.Vhost)
	r.consumerConnection, err = amqp.Dial(amqpInfo)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err)
	}

	r.producerConnection, err = amqp.Dial(amqpInfo)
	if err != nil {
		log.Fatalf("failed to connect to RabbitMQ: %s", err)
	}

	r.consumerChannel, err = r.consumerConnection.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err)
	}

	r.producerChannel, err = r.producerConnection.Channel()
	if err != nil {
		log.Fatalf("failed to open a channel: %s", err)
	}
}

// Stop stops the rabbitmq service
func (r *Rabbitmq) Stop() error {
	if !r.started {
		return nil
	}
	r.started = false
	return nil
}

// Config gets the service configuration
func (r *Rabbitmq) Config() *rabbitmqConfig.Config {
	return r.config
}

// ConfigFile gets the configuration file
func (r *Rabbitmq) ConfigFile() string {
	return configFile
}

// Produce produces to the rabbitmq
func (r *Rabbitmq) Produce(message any, exchange, routingKey string) error {
	bytesMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}
	err = r.producerChannel.Publish(
		exchange,   // Name of the exchange
		routingKey, // Routing key
		false,      // Mandatory: Return message if it cannot be routed
		false,      // Immediate: Return message if it cannot be delivered immediately
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytesMessage,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to exchange %s: %s", exchange, err)
	}

	return nil
}

// Consume consumes from the rabbitmq
func (r *Rabbitmq) Consume(app domain.IApp, queue string, handlers map[string]func(amqp.Delivery) bool) {
	messages, err := r.consumerChannel.Consume(
		queue, // Name of the queue to consume from
		fmt.Sprintf("%s-%s", r.config.Host, queue), // Consumer name
		false, // Auto-acknowledge: Remove messages from the queue once consumed
		false, // Exclusive: Queue can be accessed by multiple consumers
		false, // No-local: Do not receive messages published by this connection
		false, // No-wait: Do not wait for a response
		nil,   // Arguments
	)
	if err != nil {
		log.Fatalf("failed to consume messages at queue %s: %s", queue, err)
	}

	// Start a goroutine to process the received messages
	go r.handleMessages(messages, handlers)
}

// handleMessages handles the received messages
func (r *Rabbitmq) handleMessages(messages <-chan amqp.Delivery, handlers map[string]func(amqp.Delivery) bool) {
	for m := range messages {
		if r.handlerFunc(m, handlers)(m) {
			err := m.Ack(false)
			if err != nil {
				r.app.Logger().Log().Do(err)
			}
		} else {
			err := m.Reject(false)
			if err != nil {
				r.app.Logger().Log().Do(err)
			}
		}
	}
}

// handlerFunc gets the handler function
func (r *Rabbitmq) handlerFunc(message amqp.Delivery, handlers map[string]func(message amqp.Delivery) bool) func(amqp.Delivery) bool {
	for binding, handlerFunc := range handlers {
		if message.RoutingKey == binding {
			return func(amqp.Delivery) bool {
				defer r.recover(message)
				return handlerFunc(message)
			}
		}
	}

	return func(amqp.Delivery) bool {
		return false
	}
}

// recover recovers from a panic during message consumption
func (r *Rabbitmq) recover(message amqp.Delivery) {
	if re := recover(); re != nil {
		data, _ := json.Marshal(re)
		errorMessage := string(data)
		r.app.Logger().Log().Do(
			errors.ErrorInRabbitMQConsumer().Formats(message.Exchange, message.RoutingKey, string(message.Body)),
			&domain.LoggerInfo{
				SubType: domain.RabbitSubType,
				Msg:     errorMessage,
				Response: domain.Response{
					Body:       string(message.Body),
					StatusCode: http.StatusBadRequest,
				},
			})
	}
}

// WithAdditionalConfigType sets an additional config type
func (r *Rabbitmq) WithAdditionalConfigType(obj interface{}) domain.IRabbitMQ {
	r.additionalConfigType = obj
	return r
}

// Started true if started
func (r *Rabbitmq) Started() bool {
	return r.started
}

func (r *Rabbitmq) runMigration() error {

	dir := "./migrations/rabbitmq"
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	var rabbitMQConfig RabbitMQConfig
	for _, file := range files {
		var fileRabbitMQConfig RabbitMQConfig
		readDir, err := os.ReadFile(path.Join(dir, file.Name()))
		if err != nil {
			return err
		}

		err = json.Unmarshal(readDir, &fileRabbitMQConfig)
		if err != nil {
			return err
		}

		rabbitMQConfig.Exchanges = append(rabbitMQConfig.Exchanges, fileRabbitMQConfig.Exchanges...)
		rabbitMQConfig.Queues = append(rabbitMQConfig.Queues, fileRabbitMQConfig.Queues...)
		rabbitMQConfig.Bindings = append(rabbitMQConfig.Bindings, fileRabbitMQConfig.Bindings...)
	}

	for _, ex := range rabbitMQConfig.Exchanges {
		if err := r.consumerChannel.ExchangeDeclare(ex.Name, ex.Type, ex.Durable, false, false, false, nil); err != nil {
			return err
		}
	}

	// Declare queues
	for _, q := range rabbitMQConfig.Queues {
		if _, err := r.consumerChannel.QueueDeclare(q.Name, q.Durable, false, false, false, nil); err != nil {
			return err
		}
	}

	// Bind queues to exchanges
	for _, b := range rabbitMQConfig.Bindings {
		if err := r.consumerChannel.QueueBind(b.Destination, b.RoutingKey, b.Source, false, nil); err != nil {
			return err
		}
	}

	return nil
}
