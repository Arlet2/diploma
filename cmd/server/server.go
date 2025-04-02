package server

import (
	"context"
	"fmt"
	"log/slog"
	"push_diploma/api/handlers"
	"push_diploma/internal/service/push"
	push_s "push_diploma/internal/store/push"

	"github.com/jmoiron/sqlx"
	"github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:   "server",
	Flags:  flags,
	Action: run,
}

func run(c *cli.Context) error {
	// TODO: gracefully shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.String("postgres-user"), c.String("postgres-password"),
		c.String("postgres-host"), c.String("postgres-database"),
	)

	db, err := sqlx.ConnectContext(ctx, "postgres", databaseURL)
	if err != nil {
		slog.Error("error with connecting to database: " + err.Error())
		return nil
	}

	pushStore := push_s.NewStore(db)
	pushService := push.NewService(pushStore)
	resolver := handlers.NewResolver(pushService)

	resolver.Run()

	return nil
}
