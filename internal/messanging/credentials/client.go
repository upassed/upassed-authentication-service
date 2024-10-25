package credentials

import (
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/service/credentials"
	"github.com/wagslane/go-rabbitmq"
	"log/slog"
	"reflect"
	"runtime"
)

type rabbitClient struct {
	service          credentials.Service
	rabbitConnection *rabbitmq.Conn
	cfg              *config.Config
	log              *slog.Logger
}

func Initialize(service credentials.Service, rabbitConnection *rabbitmq.Conn, cfg *config.Config, log *slog.Logger) {
	op := runtime.FuncForPC(reflect.ValueOf(Initialize).Pointer()).Name()

	log = log.With(
		slog.String("op", op),
	)

	client := &rabbitClient{
		service:          service,
		rabbitConnection: rabbitConnection,
		cfg:              cfg,
		log:              log,
	}

	go func() {
		if err := InitializeCreateQueueConsumer(client); err != nil {
			log.Error("error while initializing student queue consumer", logging.Error(err))
			return
		}
	}()
}
