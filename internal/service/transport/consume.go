package transport

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

func (t *transport) Consume(
	ctx context.Context,
	deviceID string,
	ackWait time.Duration,
	consumerFunc jetstream.MessageHandler,
) (jetstream.ConsumeContext, error) {
	stream, err := t.js.CreateOrUpdateStream(ctx, t.getDeviceStreamConfig(deviceID))
	if err != nil {
		return nil, errors.New("error with getting stream: " + err.Error())
	}

	consumerName := t.createConsumerName(deviceID)
	consumer, err := stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    consumerName,
		Durable: consumerName,
		AckWait: ackWait,
	})
	if err != nil {
		return nil, errors.New("error with creating consumer: " + err.Error())
	}

	ch, err := consumer.Consume(consumerFunc)
	if err != nil {
		return nil, errors.New("error with start consuming: " + err.Error())
	}

	return ch, nil
}

func (t *transport) createConsumerName(deviceID string) string {
	return fmt.Sprintf("device-consumer-%v", deviceID)
}
