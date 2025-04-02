package push

import (
	"context"
	"errors"
	"log/slog"
	"push_diploma/internal/core"

	"github.com/google/uuid"
)

type service struct {
	pushStore core.PushStore
}

func NewService(
	pushStore core.PushStore,
) core.PushService {
	return &service{
		pushStore: pushStore,
	}
}

func (s *service) SendPush(
	ctx context.Context,
	push core.Push,
) (id uuid.UUID, err error) {
	push.ID = uuid.New()
	push.Status = core.PushStatusOnDelivery

	// TODO: sending to NATS

	err = s.pushStore.Create(ctx, push)
	if err != nil {
		err = errors.New("error with creating push: " + err.Error())
		slog.Error(err.Error())
		return uuid.Nil, err
	}

	return push.ID, nil
}
