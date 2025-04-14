package transport

import (
	"context"
	"errors"
	"push_diploma/internal/core"

	"github.com/nats-io/nats.go/jetstream"
)

func (t *transport) SendPush(ctx context.Context, push core.Push) (seqNumber uint64, err error) {
	_, err = t.js.CreateOrUpdateStream(ctx, t.getDeviceStreamConfig(push.DeviceID))
	if err != nil {
		return 0, errors.New("error with checking stream: " + err.Error())
	}

	msg, err := t.mapPushToProtobuf(push)
	if err != nil {
		return 0, errors.New("error with mapping to protobuf: " + err.Error())
	}

	ack, err := t.js.Publish(ctx,
		t.createDeviceTopic(push.DeviceID),
		msg,
		jetstream.WithMsgID(push.ID.String()), jetstream.WithMsgTTL(t.cfg.MessageTTL))
	if err != nil {
		return 0, errors.New("error with publishing push: " + err.Error())
	}

	return ack.Sequence, nil
}
