package messanging

import (
	"errors"
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
)

var (
	errOpeningRabbitConnection = errors.New("unable to create connection to rabbit")
)

func OpenRabbitConnection(cfg *config.Config, log *slog.Logger) (*rabbitmq.Conn, error) {
	log = logging.Wrap(log, logging.WithOp(OpenRabbitConnection))

	log.Info("creating new rabbitmq connection")
	rabbitConnection, err := rabbitmq.NewConn(
		cfg.GetRabbitConnectionString(),
		rabbitmq.WithConnectionOptionsLogging,
	)

	if err != nil {
		log.Error("unable to open connection to rabbitmq", logging.Error(err))
		return nil, errOpeningRabbitConnection
	}

	log.Info("rabbit connection established successfully")
	return rabbitConnection, nil
}
