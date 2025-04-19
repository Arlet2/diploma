package device

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"push_diploma/internal/core"
	"push_diploma/pkg/protocol"

	"github.com/google/uuid"
)

func (c *connection) handleDeliveryAck(ctx context.Context, deliveryAck *protocol.DeliveryAck) error {
	msg, ok := c.msgs[deliveryAck.PushId]
	if !ok {
		return errors.New("error with getting push: push not found")
	}

	pushID, err := uuid.Parse(deliveryAck.PushId)
	if err != nil {
		return errors.New("error with parsing push id: " + err.Error())
	}

	status := c.mapStatus(deliveryAck.Status)

	err = c.transport.SendStatusUpdate(ctx, core.StatusUpdateMsg{
		PushID:    pushID,
		NewStatus: status,
	})
	if err != nil {
		return errors.New("error with sending push status: " + err.Error())
	}

	err = msg.transportMsg.Ack()
	if err != nil {
		return errors.New("error with message ack: " + err.Error())
	}

	c.msgMutex.Lock()
	delete(c.msgs, deliveryAck.PushId)
	c.msgMutex.Unlock()

	slog.Info(fmt.Sprintf("message %v was acked", deliveryAck.PushId))

	return nil
}

func (c *connection) mapStatus(status protocol.DeliveryAck_Status) core.PushStatus {
	switch status {
	case protocol.DeliveryAck_STATUS_ACK:
		return core.PushStatusDelivered
	case protocol.DeliveryAck_STATUS_NACK:
		return core.PushStatusNacked
	default:
		slog.Warn("unknown push status: " + status.String())
		return core.PushStatusDelivered
	}
}
