package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"push_diploma/api/handlers"
	nats_r "push_diploma/api/nats"
	"push_diploma/internal/service/push"
	"push_diploma/internal/service/transport"
	push_s "push_diploma/internal/store/push"
	"syscall"

	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:   "server",
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

	postgresHost := c.String("postgres-host")
	if postgresHost == "" {
		slog.Error("postgres host cannot be empty")

		return nil
	}

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.String("postgres-user"), c.String("postgres-password"),
		postgresHost, c.String("postgres-database"),
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", databaseURL)
	if err != nil {
		slog.Error("error with connecting to database: " + err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		slog.Error(fmt.Sprintf("error with connecting to database (host: %s): "+err.Error(), postgresHost))
		return nil
	}

	slog.Info("successfully connect to Postgres at " + postgresHost)

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

	transport := transport.NewService(js, transport.Config{
		MessageTTL: c.Duration("message-ttl"),
	})

	pushStore := push_s.NewStore(db)
	pushService := push.NewService(pushStore, transport)
	resolver := handlers.NewResolver(pushService, handlers.Config{ListenHost: c.String("server-host")})

	natsResolver := nats_r.NewResolver(js, pushService)

	go resolver.Run()
	go natsResolver.Run(ctx)

	<-ctx.Done()
	resolver.Shutdown()

	return nil
}
