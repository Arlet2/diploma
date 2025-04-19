package core

import "context"

type AuthGateway interface {
	VerifyToken(ctx context.Context, token string) (deviceID string, err error)
}
