package handlers

import (
	"log/slog"
	"push_diploma/internal/core"

	"github.com/gofiber/fiber/v2"
)

type (
	Resolver struct {
		pushService core.PushService
		app         *fiber.App
		cfg         Config
	}

	Config struct {
		ListenHost string
	}
)

var (
	path = "/pushes/api/v1"
)

func NewResolver(
	pushService core.PushService,
	cfg Config,
) Resolver {
	return Resolver{
		pushService: pushService,
		app:         fiber.New(fiber.Config{}),
		cfg:         cfg,
	}
}

func (r *Resolver) Run() {
	r.app.Post(path+"/send", r.send)

	slog.Info("http server started at " + r.cfg.ListenHost)
	err := r.app.Listen(r.cfg.ListenHost)
	if err != nil {
		slog.Error("error with listening: " + err.Error())
	}
}

func (r *Resolver) Shutdown() {
	err := r.app.Shutdown()
	if err != nil {
		slog.Error("error with shutdown: " + err.Error())
	}
}
