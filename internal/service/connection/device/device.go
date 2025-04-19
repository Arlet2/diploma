package device

import (
	"context"
	"push_diploma/internal/core"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go/jetstream"
)

type (
	connection struct {
		websocketConn *websocket.Conn
		deviceID      string
		cfg           Config
		writeChan     chan []byte
		readChan      chan []byte
		ctxCancel     context.CancelFunc
		ctx           context.Context

		transport core.Transport

		msgs     map[string]deliveryMessage
		msgMutex sync.Mutex
	}

	deliveryMessage struct {
		transportMsg  jetstream.Msg
		deliveryCount uint64
	}

	Config struct {
		PingPeriod      time.Duration
		CloseWait       time.Duration
		AckWait         time.Duration
		RedeliveryCount uint64
	}
)

func NewConnection(
	websocketConn *websocket.Conn,
	deviceID string,
	cfg Config,
	transport core.Transport,
) core.DeviceConnection {
	return &connection{
		websocketConn: websocketConn,
		deviceID:      deviceID,
		cfg:           cfg,
		transport:     transport,
		msgs:          map[string]deliveryMessage{},
		msgMutex:      sync.Mutex{},

		readChan:  make(chan []byte),
		writeChan: make(chan []byte),
	}
}
