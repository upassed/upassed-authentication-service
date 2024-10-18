package app

import (
	"log/slog"

	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/server"
)

type App struct {
	Server *server.AppServer
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	const op = "app.New()"
	log = log.With(
		slog.String("op", op),
	)

	server := server.New(server.AppServerCreateParams{
		Config: config,
		Log:    log,
	})

	log.Info("app successfully created")
	return &App{
		Server: server,
	}, nil
}
