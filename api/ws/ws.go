package ws

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"push_diploma/internal/core"
	"strings"

	"github.com/gorilla/websocket"
)

type (
	Resolver struct {
		cfg               Config
		server            *http.ServeMux
		upgrader          websocket.Upgrader
		connectionService core.ConnectionService
		authGateway       core.AuthGateway
	}

	Config struct {
		Host string
	}
)

const (
	path = "/pushes/ws/v1"
)

func NewResolver(
	cfg Config,
	connectionService core.ConnectionService,
	authGateway core.AuthGateway,
) Resolver {
	return Resolver{
		cfg:               cfg,
		server:            http.NewServeMux(),
		upgrader:          websocket.Upgrader{},
		connectionService: connectionService,
		authGateway:       authGateway,
	}
}

func (r *Resolver) Run(ctx context.Context) {
	r.server.HandleFunc(path, r.handleWS)

	slog.Info("http server started at " + r.cfg.Host)
	err := http.ListenAndServe(r.cfg.Host, r.server)
	if err != nil {
		slog.Error("error with start http server: " + err.Error())
		os.Exit(1)
	}
}

func (r *Resolver) handleWS(w http.ResponseWriter, req *http.Request) {
	reqCtx := context.Background()

	authToken, err := r.getTokenFromHeader(req)
	if err != nil {
		slog.Error("error with getting token from header: " + err.Error())

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error())) // nolint
		return
	}

	deviceID, err := r.authGateway.VerifyToken(reqCtx, authToken)
	if err != nil {
		slog.Error("error with verifying token: " + err.Error())

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error with verifying token: " + err.Error())) // nolint
		return
	}

	conn, err := r.upgrader.Upgrade(w, req, nil)
	if err != nil {
		slog.Error("error with upgrade websocket: " + err.Error())

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("error with upgrade websocket: " + err.Error())) // nolint
		return
	}

	err = r.connectionService.DialWithDevice(reqCtx, conn, deviceID)
	if err != nil {
		slog.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (r *Resolver) getTokenFromHeader(req *http.Request) (string, error) {
	parts := strings.Split(req.Header.Get("Authorization"), "Bearer ")

	if len(parts) < 2 {
		return "", errors.New("bearer auth token not found")
	}

	return parts[1], nil
}
