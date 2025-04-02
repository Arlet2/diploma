package migrate

import "github.com/urfave/cli/v2"

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
}
