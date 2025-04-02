package handlers

import (
	"log/slog"
	"push_diploma/internal/core"

	"github.com/gofiber/fiber/v2"
)

type Resolver struct {
	pushService core.PushService
}

var (
	path = "/pushes/api/v1"
)

func NewResolver(
	pushService core.PushService,
) Resolver {
	return Resolver{
		pushService: pushService,
	}
}

func (r *Resolver) Run() {
	app := fiber.New(fiber.Config{})

	app.Post(path+"/send", r.send)

	// TODO: normal config
	err := app.Listen(":8080")
	if err != nil {
		slog.Error("error with listening: " + err.Error())
	}
}
