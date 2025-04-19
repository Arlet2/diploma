package core

import (
	"context"

	"github.com/gorilla/websocket"
)

type ConnectionService interface {
	DialWithDevice(ctx context.Context, websocketConn *websocket.Conn, deviceID string) error
}

type DeviceConnection interface {
	Start(ctx context.Context) error
	Handle()
}
