package auth

import (
	"context"
	"log/slog"
	"push_diploma/internal/core"
	"push_diploma/pkg/verifier"

	"google.golang.org/grpc"
)

type gateway struct {
	verifierClient verifier.VerifierClient
}

func NewGateway(
	grpcConn *grpc.ClientConn,
) core.AuthGateway {
	return &gateway{
		verifierClient: verifier.NewVerifierClient(grpcConn),
	}
}

func (g *gateway) VerifyToken(ctx context.Context, token string) (deviceID string, err error) {
	resp, err := g.verifierClient.VerifyToken(ctx, &verifier.VerifyTokenRequest{
		Token: token,
	})
	if err != nil {
		return "", core.ErrAuthServerError
	}

	switch resp.Status {
	case verifier.VerifyTokenResponse_STATUS_OK:
		return resp.DeviceId, nil
	case verifier.VerifyTokenResponse_STATUS_BAD_TOKEN:
		return "", core.ErrBadToken
	case verifier.VerifyTokenResponse_STATUS_DEVICE_NOT_FOUND:
		return "", core.ErrDeviceNotFound
	case verifier.VerifyTokenResponse_STATUS_TOKEN_EXPIRED:
		return "", core.ErrTokenExpired
	case verifier.VerifyTokenResponse_STATUS_INTERNAL_SERVER_ERROR:
		return "", core.ErrAuthServerError
	default:
		slog.Error("unknown status from auth server: " + resp.Status.String())
		return "", core.ErrAuthServerError
	}
}
