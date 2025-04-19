package nats

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"push_diploma/internal/core"

	"github.com/nats-io/nats.go/jetstream"
)

type Resolver struct {
	js          jetstream.JetStream
	pushService core.PushService
}

const (
	streamName   = "listener"
	topic        = "listener-topic"
	consumerName = "listener-consumer"
)

func NewResolver(
	js jetstream.JetStream,
	pushService core.PushService,
) *Resolver {
	return &Resolver{
		js:          js,
		pushService: pushService,
	}
}

func (r *Resolver) Run(ctx context.Context) {

	stream, err := r.createListenerStream(ctx)
	if err != nil {
		slog.Error("error with creating/updating stream: " + err.Error())
		os.Exit(1)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Name:    consumerName,
		Durable: consumerName,
	})
	if err != nil {
		slog.Error("error with creating/updating consumer: " + err.Error())
		os.Exit(1)
	}

	ch, err := consumer.Consume(r.handleMessage)
	if err != nil {
		slog.Error("error with consuming messages: " + err.Error())
		os.Exit(1)
	}

	slog.Info("nats listener is started")
	<-ctx.Done()
	ch.Stop()
}

func (r *Resolver) handleMessage(msg jetstream.Msg) {
	ctx := context.Background()

	var statusUpdateMsg core.StatusUpdateMsg
	err := json.Unmarshal(msg.Data(), &statusUpdateMsg)
	if err != nil {
		slog.Error("error with unmarshalling msg: " + err.Error())
		return
	}

	err = r.pushService.UpdateStatus(ctx, statusUpdateMsg.PushID, statusUpdateMsg.NewStatus)
	if err != nil {
		slog.Error("error with updating status: " + err.Error())
		return
	}

	err = msg.Ack()
	if err != nil {
		slog.Error("error with msg ack: " + err.Error())
		return
	}
}

func (r *Resolver) createListenerStream(ctx context.Context) (jetstream.Stream, error) {
	stream, err := r.js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{topic},
	})
	if err != nil {
		return nil, err
	}

	return stream, nil
}
