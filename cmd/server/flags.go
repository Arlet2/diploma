package server

import (
	"time"

	"github.com/urfave/cli/v2"
)

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "postgres-host",
		EnvVars: []string{"POSTGRES_HOST"},
		Value:   "localhost:5432",
	},
	&cli.StringFlag{
		Name:    "postgres-database",
		EnvVars: []string{"POSTGRES_DATABASE"},
		Value:   "diploma",
	},
	&cli.StringFlag{
		Name:    "postgres-user",
		EnvVars: []string{"POSTGRES_USER"},
		Value:   "postgres",
	},
	&cli.StringFlag{
		Name:    "postgres-password",
		EnvVars: []string{"POSTGRES_PASSWORD"},
		Value:   "mysecretpassword",
	},
	&cli.DurationFlag{
		Name:    "message-ttl",
		EnvVars: []string{"MESSAGE_TTL"},
		Value:   1 * time.Hour,
	},
	&cli.StringFlag{
		Name:    "server-host",
		EnvVars: []string{"SERVER_HOST"},
		Value:   "0.0.0.0:8080",
	},
	&cli.StringFlag{
		Name:    "nats-host",
		EnvVars: []string{"NATS_HOST"},
		Value:   "localhost:4222",
	},
}
