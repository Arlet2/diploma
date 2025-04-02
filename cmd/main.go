package main

import (
	"log/slog"
	"os"
	"push_diploma/cmd/migrate"
	"push_diploma/cmd/server"

	"github.com/urfave/cli/v2"
)

var app = cli.App{
	Name: "Diploma Push",
	Authors: []*cli.Author{
		{Name: "Artem Shulga P34111", Email: "artemshulga03@gmail.com"},
	},
	Commands: []*cli.Command{
		&migrate.Cmd,
		&server.Cmd,
	},
}

func main() {
	err := app.Run(os.Args)
	if err != nil {
		slog.Error(err.Error())
	}
}
