package core

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PushStatus string

const (
	PushStatusOnDelivery = PushStatus("ON_DELIVERY")
	PushStatusDelivered  = PushStatus("DELIVERED")
	PushStatusNacked     = PushStatus("NACKED")
)

// TODO: validation, add validation for empty strings
type Push struct {
	ID        uuid.UUID  `db:"uuid"`
	Title     string     `db:"title" validate:"required"`
	Text      string     `db:"text" validate:"required"`
	Status    PushStatus `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeviceID  string     `db:"device_id" validate:"required"`
}

type PushService interface {
	SendPush(ctx context.Context, push Push) (id uuid.UUID, err error)
}

type PushStore interface {
	Create(ctx context.Context, push Push) error
	UpdateStatus(ctx context.Context, id uuid.UUID, newStatus PushStatus) error
}
