package writer

import (
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

// Rabbitmq
type Rabbit struct {
	// client
	client domain.IRabbitMQ
	// path
	path string
	// queue
	queue string
	// routing Key
	routingKey string
	// fallback
	fallback Fallback
}

// NewRabbit creates a new writer to rabbitmq
func NewRabbit(client domain.IRabbitMQ, path, queue, routingKey string) *Rabbit {
	rabbit := &Rabbit{
		client:     client,
		path:       path,
		queue:      queue,
		routingKey: routingKey,
	}

	rabbit.fallback = Fallback{
		reader: NewFallbackReader(path, rabbit.getFallbackFileName()),
		writer: NewFallbackWriter(path, rabbit.getFallbackFileName()),
	}

	if err := rabbit.dispatchFallbackMessages(); err != nil {
		message.ErrorMessage("rabbitmq logger fallback writer", errors.ErrorRabbitmqFallback().Formats(err))
	}

	return rabbit
}

// Write writes the bytes to the writer
func (r *Rabbit) Write(message []byte) (n int, err error) {
	if r.client != nil {
		err = r.produceMessage(string(message))
	}

	if r.client == nil || err != nil {
		if n, err = r.fallback.writer.Write(message); err != nil {
			return 0, err
		}
		return n, nil
	}

	return len(message), nil
}

// dispatchFallbackMessages dispatches the fallback messages
func (r *Rabbit) dispatchFallbackMessages() (err error) {
	if r.client == nil {
		return errors.ErrorInRabbitmqClientNotFound().Formats("logging (writer)")
	}

	var lines []string
	if lines, err = r.fallback.reader.ReadLines(); err != nil {
		return err
	}

	for _, line := range lines {
		if err = r.produceMessage(line); err != nil {
			return err
		}
	}

	return r.fallback.writer.Remove()
}

// produceMessage produces the messages to the rabbitmq
func (r *Rabbit) produceMessage(message string) error {
	return r.client.Produce(message, r.queue, r.routingKey)
}

// getFallbackFileName gets the error file name
func (r *Rabbit) getFallbackFileName() string {
	return "fallback_rabbitmq.log"
}
