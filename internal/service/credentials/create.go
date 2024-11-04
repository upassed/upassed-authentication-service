package credentials

import (
	"context"
	"errors"
	"github.com/upassed/upassed-authentication-service/internal/async"
	"github.com/upassed/upassed-authentication-service/internal/handling"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"log/slog"
)

var (
	ErrCreateCredentialsDeadlineExceeded = errors.New("create credentials deadline exceeded")
)

func (service *credentialsServiceImpl) Create(ctx context.Context, credentials *business.Credentials) (*business.CreateCredentialsResponse, error) {
	log := logging.Wrap(service.log,
		logging.WithOp(service.Create),
		logging.WithCtx(ctx),
		logging.WithAny("username", credentials.Username),
	)

	spanContext, span := otel.Tracer(service.cfg.Tracing.CredentialsTracerName).Start(ctx, "credentialsService#Create")
	span.SetAttributes(attribute.String("username", credentials.Username))
	defer span.End()

	log.Info("started creating credentials")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	credentialsCreateResponse, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.CreateCredentialsResponse, error) {
		log.Info("checking credentials duplicate exists")
		duplicateExists, err := service.repository.CheckDuplicatesExists(ctx, credentials.Username)
		if err != nil {
			log.Error("error while checking credentials duplicates", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, err
		}

		if duplicateExists {
			log.Error("credentials with this username already exists")
			tracing.SetSpanError(span, errors.New("credentials duplicate found"))
			return nil, handling.Wrap(errors.New("credentials duplicate found"), handling.WithCode(codes.AlreadyExists))
		}

		log.Info("converting to domain credentials, generating password hash")
		domainCredentials, err := ConvertToDomainCredentials(credentials)
		if err != nil {
			log.Error("unable to convert to domain credentials", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, handling.Wrap(errors.New("error generating password hash"), handling.WithCode(codes.Internal))
		}

		log.Info("saving credentials to the database")
		if err := service.repository.Save(ctx, domainCredentials); err != nil {
			log.Error("error while saving credentials to a database", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, handling.Process(err)
		}

		log.Info("credentials successfully created")
		return &business.CreateCredentialsResponse{
			CreatedCredentialsID: domainCredentials.ID,
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("credentials creating deadline exceeded")
			tracing.SetSpanError(span, err)
			return nil, handling.Wrap(ErrCreateCredentialsDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while creating credentials", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.Process(err)
	}

	log.Info("credentials successfully created", slog.Any("createdCredentialsID", credentialsCreateResponse.CreatedCredentialsID))
	return credentialsCreateResponse, nil
}
