package transport

import (
	"fmt"
	"push_diploma/internal/core"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type (
	transport struct {
		js  jetstream.JetStream
		cfg Config
	}

	Config struct {
		MessageTTL time.Duration
	}
)

func NewService(
	js jetstream.JetStream,
	cfg Config,
) core.Transport {
	return &transport{
		js:  js,
		cfg: cfg,
	}
}

func (t *transport) getDeviceStreamConfig(deviceID string) jetstream.StreamConfig {
	return jetstream.StreamConfig{
		Name:        t.createDeviceStreamName(deviceID),
		Subjects:    []string{t.createDeviceTopic(deviceID)},
		Retention:   jetstream.WorkQueuePolicy,
		AllowMsgTTL: true,
		AllowDirect: true,
	}
}

func (t *transport) createDeviceStreamName(deviceID string) string {
	return fmt.Sprintf("device-stream-%s", deviceID)
}

func (t *transport) createDeviceTopic(deviceID string) string {
	return fmt.Sprintf("device-topic-%s", deviceID)
}
