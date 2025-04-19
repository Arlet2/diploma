package push

import (
	"context"
	"push_diploma/internal/core"

	"github.com/google/uuid"
)

func (s *service) UpdateStatus(ctx context.Context, pushID uuid.UUID, newStatus core.PushStatus) error {
	err := s.pushStore.UpdateStatus(ctx, pushID, newStatus)
	if err != nil {
		return err
	}

	return nil
}
