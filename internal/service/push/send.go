package push

import (
	"context"
	"errors"
	"log/slog"
	"push_diploma/internal/core"
	"time"

	"github.com/google/uuid"
)

func (s *service) SendPush(
	ctx context.Context,
	push core.Push,
) (id uuid.UUID, err error) {
	push.ID = uuid.New()
	push.Status = core.PushStatusOnDelivery
	push.CreatedAt = time.Now()

	err = s.validator.Struct(&push)
	if err != nil {
		slog.Error("validation error: " + err.Error())

		return uuid.Nil, err
	}

	seqNumber, err := s.transport.SendPush(ctx, push)
	if err != nil {
		err := errors.New("error with sending push: " + err.Error())
		slog.Error(err.Error())

		return uuid.Nil, err
	}

	err = s.pushStore.Create(ctx, push)
	if err != nil {
		transportErr := s.transport.CancelPush(ctx, push.DeviceID, seqNumber)
		if transportErr != nil {
			slog.Error("error with cancelling push: " + err.Error())
		}

		err = errors.New("error with creating push: " + err.Error())
		slog.Error(err.Error())
		return uuid.Nil, err
	}

	return push.ID, nil
}
