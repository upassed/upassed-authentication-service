package credentials

import (
	"context"
	"github.com/google/uuid"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/middleware"
	"github.com/wagslane/go-rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
	"reflect"
	"runtime"
)

func (client *rabbitClient) CreateQueueConsumer(log *slog.Logger) func(d rabbitmq.Delivery) rabbitmq.Action {
	op := runtime.FuncForPC(reflect.ValueOf(client.CreateQueueConsumer).Pointer()).Name()

	return func(delivery rabbitmq.Delivery) rabbitmq.Action {
		requestID := uuid.New().String()
		ctx := context.WithValue(context.Background(), middleware.RequestIDKey, requestID)

		log = log.With(
			slog.String("op", op),
			slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
		)

		log.Info("consumed credentials create message", slog.String("messageBody", string(delivery.Body)))
		spanContext, span := otel.Tracer(client.cfg.Tracing.CredentialsTracerName).Start(ctx, "credentials#Create")
		span.SetAttributes(attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)))
		defer span.End()

		request, err := ConvertToCredentialsCreateRequest(delivery.Body)
		if err != nil {
			log.Error("error parsing message body to json", logging.Error(err))
			return rabbitmq.NackDiscard
		}

		if err := request.Validate(); err != nil {
			log.Error("request is invalid", logging.Error(err))
			return rabbitmq.NackDiscard
		}

		response, err := client.service.Create(spanContext, ConvertToCredentials(request))
		if err != nil {
			log.Error("error while creating credentials", logging.Error(err))
			return rabbitmq.NackDiscard
		}

		log.Info("successfully created credentials", slog.Any("createdCredentialsID", response.CreatedCredentialsID))
		return rabbitmq.Ack
	}
}
