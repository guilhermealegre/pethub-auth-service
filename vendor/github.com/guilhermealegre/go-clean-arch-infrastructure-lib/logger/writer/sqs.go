package writer

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/domain/message"
	"github.com/guilhermealegre/go-clean-arch-infrastructure-lib/errors"
)

// SQS
type SQS struct {
	// client
	client domain.ISQSConnection
	// path
	path string
	// queue
	queue string
	// message attributes
	messageAttributes map[string]*sqs.MessageAttributeValue
	// fallback
	fallback Fallback
}

// NewSQS creates a new writer to SQS
func NewSQS(client domain.ISQS, path, connection, queue, routingKey string) *SQS {
	messageAttributes := make(map[string]*sqs.MessageAttributeValue)
	messageAttributes["handler"] = &sqs.MessageAttributeValue{
		DataType:    aws.String("String"),
		StringValue: aws.String(routingKey),
	}

	sqs := &SQS{
		client:            client.Connection(connection),
		path:              path,
		queue:             queue,
		messageAttributes: messageAttributes,
	}

	sqs.fallback = Fallback{
		reader: NewFallbackReader(path, sqs.getFallbackFileName()),
		writer: NewFallbackWriter(path, sqs.getFallbackFileName()),
	}

	if err := sqs.dispatchFallbackMessages(); err != nil {
		message.ErrorMessage("sqs logger fallback writer", errors.ErrorSQSFallback().Formats(err))
	}

	return sqs
}

// Write writes the bytes to the writer
func (r *SQS) Write(message []byte) (n int, err error) {
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
func (r *SQS) dispatchFallbackMessages() error {
	if r.client == nil {
		return errors.ErrorInSQSClientNotFound().Formats("logging (writer)")
	}

	var lines []string
	var err error
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

// produceMessage produces the messages to the SQS
func (r *SQS) produceMessage(message string) error {
	return r.client.Produce(context.Background(), r.queue, r.messageAttributes, message)
}

// getFallbackFileName gets the error file name
func (r *SQS) getFallbackFileName() string {
	return "fallback_sqs.log"
}
