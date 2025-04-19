package transport

import (
	"context"
	"encoding/json"
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

func (t *transport) SendStatusUpdate(ctx context.Context, statusUpdate core.StatusUpdateMsg) error {
	data, err := json.Marshal(statusUpdate)
	if err != nil {
		return errors.New("error with marshalling json: " + err.Error())
	}

	_, err = t.js.Publish(ctx,
		t.createListenerTopic(),
		data,
	)
	if err != nil {
		return errors.New("error with publishing message to nats: " + err.Error())
	}

	return nil
}
