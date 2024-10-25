package app

import (
	"log/slog"
	"reflect"
	"runtime"

	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/server"
)

type App struct {
	Server *server.AppServer
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	op := runtime.FuncForPC(reflect.ValueOf(New).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	appServer := server.New(server.AppServerCreateParams{
		Config: config,
		Log:    log,
	})

	log.Info("app successfully created")
	return &App{
		Server: appServer,
	}, nil
}
