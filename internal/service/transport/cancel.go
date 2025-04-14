package transport

import (
	"context"
	"errors"
)

func (t *transport) CancelPush(ctx context.Context, deviceID string, seqNumber uint64) error {
	stream, err := t.js.CreateOrUpdateStream(ctx, t.getDeviceStreamConfig(deviceID))
	if err != nil {
		return errors.New("error with getting stream: " + err.Error())
	}

	err = stream.DeleteMsg(ctx, seqNumber)
	if err != nil {
		return errors.New("error with deleting message: " + err.Error())
	}

	return nil
}
