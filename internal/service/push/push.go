package push

import (
	"push_diploma/internal/core"

	"github.com/go-playground/validator/v10"
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
