package ws

import "github.com/urfave/cli/v2"

var flags = []cli.Flag{
	&cli.StringFlag{
		Name:    "ws-host",
		EnvVars: []string{"WS_HOST"},
		Value:   "0.0.0.0:9000",
	},
	&cli.StringFlag{
		Name:    "nats-host",
		EnvVars: []string{"NATS_HOST"},
		Value:   "localhost:4222",
	},
	&cli.StringFlag{
		Name:    "grpc-auth-host",
		EnvVars: []string{"GRPC_AUTH_HOST"},
		Value:   "localhost:9090",
	},
}
