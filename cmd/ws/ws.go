package ws

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"push_diploma/api/ws"
	"push_diploma/internal/gateway/auth"
	"push_diploma/internal/service/connection"
	"push_diploma/internal/service/connection/device"
	"push_diploma/internal/service/transport"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var Cmd = cli.Command{
	Name:   "ws",
	Flags:  flags,
	Action: run,
}

func run(c *cli.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
			return
		case s := <-sig:
			slog.Info(fmt.Sprintf("signal %s received", s.String()))
			cancel()
		}
	}()

	resolverCfg := ws.Config{
		Host: c.String("ws-host"),
	}

	deviceCfg := device.Config{
		PingPeriod:      10 * time.Second,
		CloseWait:       5 * time.Second,
		AckWait:         1 * time.Second,
		RedeliveryCount: 2,
	}

	natsHost := c.String("nats-host")
	if natsHost == "" {
		slog.Error("nats host cannot be empty")

		return nil
	}

	natsConn, err := nats.Connect(natsHost, nats.DisconnectErrHandler(func(c *nats.Conn, err error) {
		slog.Error("nats was disconnected: " + err.Error())

		os.Exit(1)
	}))
	if err != nil {
		slog.Error("error with connecting to nats: " + err.Error())
		return nil
	}

	js, err := jetstream.New(natsConn)
	if err != nil {
		slog.Error("error with connecting to nats as jetstream: " + err.Error())
		return nil
	}

	slog.Info("successfully connect to nats jetstream at " + natsHost)

	transportService := transport.NewService(js, transport.Config{
		MessageTTL: c.Duration("message-ttl"),
	})

	connectionService := connection.NewConnectionService(transportService, deviceCfg)

	grpcAuthHost := c.String("grpc-auth-host")
	if grpcAuthHost == "" {
		return errors.New("grpc auth host cannot be empty")
	}

	grpcConn, err := grpc.NewClient(grpcAuthHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New("error with connecting to auth service by grpc: " + err.Error())
	}

	authGateway := auth.NewGateway(grpcConn)

	resolver := ws.NewResolver(resolverCfg, connectionService, authGateway)

	go resolver.Run(ctx)

	<-ctx.Done()

	return nil
}
