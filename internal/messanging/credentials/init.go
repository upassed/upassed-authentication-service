package credentials

import (
	"errors"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/wagslane/go-rabbitmq"
)

var (
	errCreatingCredentialsCreateQueueConsumer = errors.New("unable to create credentials queue consumer")
	errRunningCredentialsCreateQueueConsumer  = errors.New("unable to run credentials queue consumer")
)

func InitializeCreateQueueConsumer(client *rabbitClient) error {
	log := logging.Wrap(client.log, logging.WithOp(InitializeCreateQueueConsumer))

	log.Info("started creating credentials create queue consumer")
	credentialsCreateQueueConsumer, err := rabbitmq.NewConsumer(
		client.rabbitConnection,
		client.cfg.Rabbit.Queues.CredentialsCreate.Name,
		rabbitmq.WithConsumerOptionsRoutingKey(client.cfg.Rabbit.Queues.CredentialsCreate.RoutingKey),
		rabbitmq.WithConsumerOptionsExchangeName(client.cfg.Rabbit.Exchange.Name),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)

	if err != nil {
		log.Error("unable to create credentials queue consumer", logging.Error(err))
		return errCreatingCredentialsCreateQueueConsumer
	}

	defer credentialsCreateQueueConsumer.Close()
	if err := credentialsCreateQueueConsumer.Run(client.CreateQueueConsumer(log)); err != nil {
		log.Error("unable to run credentials queue consumer")
		return errRunningCredentialsCreateQueueConsumer
	}

	log.Info("credentials queue consumer initialized successfully")
	return nil
}
