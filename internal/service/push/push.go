package push

import (
	"context"
	"errors"
	"log/slog"
	"push_diploma/internal/core"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type service struct {
	pushStore core.PushStore
	validator *validator.Validate
	transport core.Transport
}

func NewService(
	pushStore core.PushStore,
	transport core.Transport,
) core.PushService {
	return &service{
		pushStore: pushStore,
		validator: validator.New(),
		transport: transport,
	}
}

func (s *service) SendPush(
	ctx context.Context,
	push core.Push,
) (id uuid.UUID, err error) {
	push.ID = uuid.New()
	push.Status = core.PushStatusOnDelivery

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
