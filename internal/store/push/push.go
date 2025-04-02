package push

import (
	"context"
	"errors"
	"push_diploma/internal/core"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type store struct {
	postgresClient *sqlx.DB
}

func NewStore(
	postgresClient *sqlx.DB,
) core.PushStore {
	return &store{
		postgresClient: postgresClient,
	}
}

func (s *store) Create(ctx context.Context, push core.Push) error {
	_, err := s.postgresClient.ExecContext(ctx,
		`INSERT INTO pushes
		(id, created_at, updated_at, device_id, title, text, status) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		push.ID, time.Now().UTC(), time.Now().UTC(), push.DeviceID, push.Title, push.Text, push.Status,
	)
	if err != nil {
		return errors.New("error with push creating: " + err.Error())
	}

	return nil
}

func (s *store) UpdateStatus(ctx context.Context, id uuid.UUID, newStatus core.PushStatus) error {
	_, err := s.postgresClient.ExecContext(ctx,
		`UPDATE pushes SET status=$1 WHERE id = $2`,
		newStatus, id,
	)
	if err != nil {
		return errors.New("error with push status updating: " + err.Error())
	}

	return nil
}
