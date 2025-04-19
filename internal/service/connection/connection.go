package connection

import (
	"context"
	"log/slog"
	"push_diploma/internal/core"
	"push_diploma/internal/service/connection/device"

	"github.com/gorilla/websocket"
)

type connectionService struct {
	transport core.Transport
	deviceCfg device.Config
}

func NewConnectionService(
	transport core.Transport,
	cfg device.Config,
) core.ConnectionService {
	return &connectionService{
		transport: transport,
		deviceCfg: cfg,
	}
}

func (c *connectionService) DialWithDevice(ctx context.Context, websocketConn *websocket.Conn, deviceID string) error {
	conn := device.NewConnection(websocketConn, deviceID, c.deviceCfg, c.transport)

	err := conn.Start(ctx)
	if err != nil {
		slog.Error("error with starting connection with device: " + err.Error())
	}
	// в будущем здесь может быть процедура рукопожатия или иная проверка клиента внутри websocket

	go conn.Handle()

	return nil
}
