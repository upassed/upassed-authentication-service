package credentials

import (
	"context"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/middleware/amqp"
	loggingMiddleware "github.com/upassed/upassed-authentication-service/internal/middleware/amqp/logging"
	"github.com/upassed/upassed-authentication-service/internal/middleware/amqp/recovery"
	requestidMiddleware "github.com/upassed/upassed-authentication-service/internal/middleware/amqp/request_id"
	requestid "github.com/upassed/upassed-authentication-service/internal/middleware/common/request_id"
	"github.com/upassed/upassed-authentication-service/internal/tracing"
	"github.com/wagslane/go-rabbitmq"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"log/slog"
)

func (client *rabbitClient) CreateQueueConsumer(log *slog.Logger) func(d rabbitmq.Delivery) rabbitmq.Action {
	baseHandler := func(ctx context.Context, delivery rabbitmq.Delivery) rabbitmq.Action {
		log = logging.Wrap(log,
			logging.WithOp(client.CreateQueueConsumer),
			logging.WithCtx(ctx),
		)

		log.Info("consumed credentials create message", slog.String("messageBody", string(delivery.Body)))
		spanContext, span := otel.Tracer(client.cfg.Tracing.CredentialsTracerName).Start(ctx, "credentials#Create")
		span.SetAttributes(attribute.String(string(requestid.ContextKey), requestid.GetRequestIDFromContext(ctx)))
		defer span.End()

		log.Info("converting message body to credentials create request")
		request, err := ConvertToCredentialsCreateRequest(delivery.Body)
		if err != nil {
			log.Error("error parsing message body to json", logging.Error(err))
			tracing.SetSpanError(span, err)
			return rabbitmq.NackDiscard
		}

		span.SetAttributes(attribute.String("username", request.Username))
		log.Info("validating the credentials create request")
		if err := request.Validate(); err != nil {
			log.Error("request is invalid", logging.Error(err))
			tracing.SetSpanError(span, err)
			return rabbitmq.NackDiscard
		}

		log.Info("creating credentials")
		response, err := client.service.Create(spanContext, ConvertToCredentials(request))
		if err != nil {
			log.Error("error while creating credentials", logging.Error(err))
			tracing.SetSpanError(span, err)
			return rabbitmq.NackDiscard
		}

		log.Info("successfully created credentials", slog.Any("createdCredentialsID", response.CreatedCredentialsID))
		return rabbitmq.Ack
	}

	handlerWithMiddleware := amqp.ChainMiddleware(
		baseHandler,
		requestidMiddleware.Middleware(),
		recovery.Middleware(client.log),
		loggingMiddleware.Middleware(client.log),
	)

	return func(d rabbitmq.Delivery) (action rabbitmq.Action) {
		ctx := context.Background()
		return handlerWithMiddleware(ctx, d)
	}
}
