package device

import (
	"context"
	"log/slog"
	"push_diploma/pkg/protocol"

	"github.com/nats-io/nats.go/jetstream"
	"google.golang.org/protobuf/proto"
)

func (c *connection) Handle() {
	defer c.ctxCancel()

	ch, err := c.transport.Consume(c.ctx, c.deviceID, c.cfg.AckWait, c.handleMessage)
	if err != nil {
		slog.Error("error with starting consuming: " + err.Error())

		return
	}

	go c.replyHandling(c.ctx)

	slog.Info("connection with " + c.deviceID + " was started")
	<-c.ctx.Done()
	ch.Stop()

	slog.Info("stop consuming")
}

func (c *connection) handleMessage(msg jetstream.Msg) {
	msgID := msg.Headers().Get(jetstream.MsgIDHeader)

	deliveryMsg, ok := c.msgs[msgID]
	if !ok {
		deliveryMsg = deliveryMessage{
			transportMsg:  msg,
			deliveryCount: 0,
		}
	}

	if deliveryMsg.deliveryCount >= c.cfg.RedeliveryCount {
		c.Close()
		return
	}

	deliveryMsg.deliveryCount++
	c.msgMutex.Lock()
	c.msgs[msgID] = deliveryMsg
	c.msgMutex.Unlock()

	c.writeChan <- msg.Data()
}

func (c *connection) replyHandling(ctx context.Context) {
	defer c.ctxCancel()

	for {
		select {
		case <-ctx.Done():
			return
		case data, ok := <-c.readChan:
			if !ok {
				slog.Info("read chan was closed")
				return
			}

			var clientMessage protocol.ClientMessage
			err := proto.Unmarshal(data, &clientMessage)
			if err != nil {
				slog.Error("error with unmarshalling client message: " + err.Error())
				continue
			}

			switch val := clientMessage.ClientMessage.(type) {
			case *protocol.ClientMessage_DeliveryAck:
				err := c.handleDeliveryAck(ctx, val.DeliveryAck)
				if err != nil {
					slog.Error("error with handling delivery ack: " + err.Error())
					continue
				}
			}
		}
	}
}
