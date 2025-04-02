package migrate

import (
	"fmt"
	"log/slog"
	"push_diploma/internal/schema"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	"github.com/urfave/cli/v2"
)

var Cmd = cli.Command{
	Name:   "migrate",
	Flags:  flags,
	Action: run,
}

func run(c *cli.Context) error {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		c.String("postgres-user"), c.String("postgres-password"),
		c.String("postgres-host"), c.String("postgres-database"),
	)

	err := schema.Migrate(databaseURL)
	if err != nil {
		return err
	}

	slog.Info("Migration applied")

	return nil
}
