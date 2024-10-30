package app

import (
	"github.com/upassed/upassed-authentication-service/internal/caching"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/jwt"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/messanging"
	credentialsRabbit "github.com/upassed/upassed-authentication-service/internal/messanging/credentials"
	"github.com/upassed/upassed-authentication-service/internal/repository"
	credentialsRepo "github.com/upassed/upassed-authentication-service/internal/repository/credentials"
	"github.com/upassed/upassed-authentication-service/internal/server"
	"github.com/upassed/upassed-authentication-service/internal/service/credentials"
	"github.com/upassed/upassed-authentication-service/internal/service/token"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
)

type App struct {
	Server     *server.AppServer
	RabbitConn *rabbitmq.Conn
}

func New(config *config.Config, log *slog.Logger) (*App, error) {
	log = logging.Wrap(log, logging.WithOp(New))
	log.Info("started initializing application")

	db, err := repository.OpenGormDbConnection(config, log)
	if err != nil {
		return nil, err
	}

	redis, err := caching.OpenRedisConnection(config, log)
	if err != nil {
		return nil, err
	}

	rabbit, err := messanging.OpenRabbitConnection(config, log)
	if err != nil {
		return nil, err
	}

	credentialsRepository := credentialsRepo.New(db, redis, config, log)

	credentialsService := credentials.New(config, log, credentialsRepository)
	credentialsRabbit.Initialize(credentialsService, rabbit, config, log)

	tokenGenerator := jwt.New(config, log)

	appServer := server.New(server.AppServerCreateParams{
		Config:       config,
		Log:          log,
		TokenService: token.New(config, log, tokenGenerator, credentialsRepository),
	})

	log.Info("app successfully created")
	return &App{
		Server:     appServer,
		RabbitConn: rabbit,
	}, nil
}
