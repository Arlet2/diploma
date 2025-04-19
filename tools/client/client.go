package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"push_diploma/pkg/protocol"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const (
	host   = "ws://localhost:9000"
	wsPath = "/pushes/ws/v1"
)

const (
	authToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDQ5NzUwNDEsImlhdCI6MTc0NDk2Nzg0MSwic3ViIjoidGVzdC1kZXYtaWQifQ.UU5ZJgNt-jg4lusBq7tjQyHgYXdyKNHdQCBNaBkSQyML92eC2aXy1lZCfJfzkB-q2-S-7zwx_O-Vt47BEoaSxQ"
)

func main() {
	ctx := context.Background()

	conn, resp, err := websocket.DefaultDialer.Dial(host+wsPath, http.Header{
		"Authorization": []string{"Bearer " + authToken},
	})
	if err != nil {
		slog.Error("error with dialing with server: " + err.Error())

		var body []byte
		_, err := resp.Body.Read(body) // nolint
		if err != nil {
			slog.Error("error with reading body: " + err.Error())
		}

		slog.Error("resp body: " + string(body))

		return
	}

	readChan := make(chan []byte)
	writeChan := make(chan []byte)

	go handleWrite(ctx, conn, writeChan)
	go handleRead(ctx, conn, readChan)

	slog.Info("connection success")
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-readChan:
			if !ok {
				slog.Info("read chan was closed")
				os.Exit(1)
			}

			var serverMessage protocol.ServerMessage
			err := proto.Unmarshal(msg, &serverMessage)
			if err != nil {
				slog.Error("error with unmarshalling msg: " + err.Error())
				continue
			}

			switch val := serverMessage.ServerMessage.(type) {
			case *protocol.ServerMessage_PushNotification:
				slog.Info(fmt.Sprintf("received push: %+v", val.PushNotification))

				err := sendAck(writeChan, val.PushNotification.Id)
				if err != nil {
					slog.Error("error with sending ack: " + err.Error())
				}
			}
		}
	}

}

func handleWrite(ctx context.Context, conn *websocket.Conn, writeChan chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-writeChan:
			if !ok {
				slog.Error("write message was closed")
				return
			}

			err := conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				slog.Error("error with writing: " + err.Error())
				os.Exit(1)
			}
		}
	}
}

func handleRead(ctx context.Context, conn *websocket.Conn, readChan chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			slog.Error("error with reading: " + err.Error())
			os.Exit(1)
		}

		readChan <- msg
	}
}

func sendAck(writeChan chan []byte, pushID string) error {
	clientMessage := protocol.ClientMessage{
		ClientMessage: &protocol.ClientMessage_DeliveryAck{
			DeliveryAck: &protocol.DeliveryAck{
				PushId: pushID,
				Status: protocol.DeliveryAck_STATUS_ACK,
			},
		},
	}

	msg, err := proto.Marshal(&clientMessage)
	if err != nil {
		return errors.New("error with marshalling client message: " + err.Error())
	}

	writeChan <- msg

	return nil
}
