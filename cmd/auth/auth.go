package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"push_diploma/pkg/verifier"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "jwt-secret",
		EnvVars: []string{"JWT_SECRET"},
		Value:   "super_secret",
	},
	&cli.StringFlag{
		Name:    "device-id",
		EnvVars: []string{"DEVICE_ID"},
	},
	&cli.StringFlag{
		Name:    "grpc-host",
		EnvVars: []string{"GRPC_HOST"},
		Value:   ":9090",
	},
}

var app = cli.App{
	Name: "Organization auth server",
	Authors: []*cli.Author{
		{Name: "Artem Shulga P34111", Email: "artemshulga03@gmail.com"},
	},
	Commands: []*cli.Command{
		&authServerCmd,
		&tokenCreateCmd,
	},
	Flags: flags,
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		slog.Error(err.Error())
	}
}

var authServerCmd = cli.Command{
	Name:   "auth-server",
	Action: authServer,
	Flags:  flags,
}

func authServer(c *cli.Context) error {
	grpcServer := grpc.NewServer()

	verifier.RegisterVerifierServer(grpcServer, NewGRPCService(c.String("jwt-secret")))

	host := c.String("grpc-host")
	if host == "" {
		return errors.New("grpc host cannot be empty")
	}

	listener, err := net.Listen("tcp", host)
	if err != nil {
		return err
	}

	slog.Info("auth server started at " + host)
	err = grpcServer.Serve(listener)
	if err != nil {
		slog.Error("error with grpc serving: " + err.Error())
	}

	return nil
}

var tokenCreateCmd = cli.Command{
	Name:  "token-create",
	Flags: flags,
	Action: func(c *cli.Context) error {
		deviceID := c.String("device-id")
		if deviceID == "" {
			return errors.New("device ID cannot be empty")
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS512,
			jwt.MapClaims{
				"sub": deviceID,
				"iat": time.Now().Unix(),
				"exp": time.Now().Add(2 * time.Hour).Unix(),
			},
		)

		signedToken, err := token.SignedString([]byte(c.String("jwt-secret")))
		if err != nil {
			return err
		}

		fmt.Printf("JWT token for device-id=%v: %v\n", deviceID, signedToken)
		return nil
	},
}

type grpcService struct {
	verifier.UnimplementedVerifierServer
	jwtSecret string
}

func NewGRPCService(jwtSecret string) *grpcService {
	return &grpcService{
		jwtSecret: jwtSecret,
	}
}

func (g *grpcService) VerifyToken(ctx context.Context, req *verifier.VerifyTokenRequest) (*verifier.VerifyTokenResponse, error) {
	var status verifier.VerifyTokenResponse_Status
	var deviceID string
	token, err := jwt.Parse(req.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(g.jwtSecret), nil
	})
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenSignatureInvalid) || errors.Is(err, jwt.ErrTokenMalformed):
			status = verifier.VerifyTokenResponse_STATUS_BAD_TOKEN
		case errors.Is(err, jwt.ErrTokenExpired):
			status = verifier.VerifyTokenResponse_STATUS_TOKEN_EXPIRED
		default:
			slog.Error("unknown error: " + err.Error())
			status = verifier.VerifyTokenResponse_STATUS_INTERNAL_SERVER_ERROR
		}

	} else {
		status = verifier.VerifyTokenResponse_STATUS_OK
		deviceID, err = token.Claims.GetSubject()
		if err != nil {
			slog.Error("sub not found in token: " + err.Error())
			return &verifier.VerifyTokenResponse{
				Status: verifier.VerifyTokenResponse_STATUS_BAD_TOKEN,
			}, nil
		}
	}

	// возможны ошибки INTERNAL и DEVICE NOT FOUND, но здесь эта логика опущена для упрощения

	return &verifier.VerifyTokenResponse{
		DeviceId: deviceID,
		Status:   status,
	}, nil
}
