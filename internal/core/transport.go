package core

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go/jetstream"
)

type StatusUpdateMsg struct {
	PushID    uuid.UUID
	NewStatus PushStatus
}

type Transport interface {
	SendPush(ctx context.Context, push Push) (seqNumber uint64, err error)
	CancelPush(ctx context.Context, deviceID string, seqNumber uint64) error
	Consume(ctx context.Context, deviceID string, ackWait time.Duration, consumerFunc jetstream.MessageHandler) (jetstream.ConsumeContext, error)
	SendStatusUpdate(ctx context.Context, statusUpdate StatusUpdateMsg) error
}
