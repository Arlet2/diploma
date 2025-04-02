package schema

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
)

//go:embed migrations/*
var migrationsFS embed.FS

func Migrate(databaseURL string) error {
	migrations, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return errors.New("error with getting migration files: " + err.Error())
	}

	m, err := migrate.NewWithSourceInstance(
		"iofs",
		migrations,
		databaseURL,
	)
	if err != nil {
		return errors.New("error with creating migrations: " + err.Error())
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return errors.New("error with migration up: " + err.Error())
	}

	return nil
}
