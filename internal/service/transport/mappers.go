package transport

import (
	"push_diploma/internal/core"
	"push_diploma/pkg/protocol"
	"time"

	"google.golang.org/protobuf/proto"
)

func (t *transport) mapPushToProtobuf(push core.Push) ([]byte, error) {
	serverMessage := protocol.ServerMessage{
		ServerMessage: &protocol.ServerMessage_PushNotification{
			PushNotification: &protocol.PushNotification{
				Id:        push.ID.String(),
				CreatedAt: push.CreatedAt.Format(time.RFC3339),
				Title:     push.Title,
				Body:      push.Text,
			},
		},
	}

	data, err := proto.Marshal(&serverMessage)
	if err != nil {
		return nil, err
	}

	return data, nil
}
