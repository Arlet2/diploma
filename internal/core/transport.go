package core

import "context"

type Transport interface {
	SendPush(ctx context.Context, push Push) (seqNumber uint64, err error)
	CancelPush(ctx context.Context, deviceID string, seqNumber uint64) error
}
