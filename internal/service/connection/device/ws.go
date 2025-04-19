package device

import (
	"context"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

func (c *connection) Write(ctx context.Context) {
	ticker := time.NewTicker(c.cfg.PingPeriod)

	defer func() {
		ticker.Stop()
		slog.Info("stop writing")
		c.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.writeChan:
			if !ok {
				slog.Info("writing channel was closed")
				return
			}

			c.websocketConn.SetWriteDeadline(time.Now().Add(2 * c.cfg.PingPeriod)) // nolint

			err := c.websocketConn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				slog.Error("error with writing message: " + err.Error())
				return
			}
		case <-ticker.C:
			err := c.websocketConn.WriteControl(websocket.PingMessage, nil, time.Now().Add(2*c.cfg.PingPeriod))
			if err != nil {
				slog.Error("error with writing ping: " + err.Error())
				return
			}
		}
	}
}

func (c *connection) Read(ctx context.Context) {
	defer func() {
		slog.Info("stop reading")
		c.Close()
		close(c.readChan)
	}()

	for {
		_, message, err := c.websocketConn.ReadMessage()
		if err != nil {
			slog.Error("error with reading message from websocket: " + err.Error())
			return
		}

		c.readChan <- message
	}
}

func (c *connection) Close() {
	c.ctxCancel()

	err := c.websocketConn.WriteControl(websocket.CloseMessage, nil, time.Now().Add(c.cfg.CloseWait))
	if err != nil {
		slog.Error("error with sending close message: " + err.Error())
		return
	}

	err = c.websocketConn.Close()
	if err != nil {
		slog.Error("error with close connection: " + err.Error())
		return
	}

	slog.Info("success connection closing")
}
