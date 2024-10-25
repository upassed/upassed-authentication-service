package credentials

import (
	"context"
	"errors"
	"github.com/upassed/upassed-authentication-service/internal/async"
	"github.com/upassed/upassed-authentication-service/internal/handling"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/middleware"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"log/slog"
	"reflect"
	"runtime"
)

var (
	errCreateCredentialsDeadlineExceeded = errors.New("create credentials deadline exceeded")
)

func (service *credentialsServiceImpl) Create(ctx context.Context, credentials business.Credentials) (business.CreateCredentialsResponse, error) {
	op := runtime.FuncForPC(reflect.ValueOf(service.Create).Pointer()).Name()

	log := service.log.With(
		slog.String("op", op),
		slog.String("username", credentials.Username),
		slog.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)

	spanContext, span := otel.Tracer(service.cfg.Tracing.CredentialsTracerName).Start(ctx, "credentialsService#Create")
	defer span.End()

	log.Info("started creating credentials")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	credentialsCreateResponse, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (business.CreateCredentialsResponse, error) {
		duplicateExists, err := service.repository.CheckDuplicatesExists(ctx, credentials.Username)
		if err != nil {
			return business.CreateCredentialsResponse{}, err
		}

		if duplicateExists {
			log.Error("credentials with this username already exists")
			return business.CreateCredentialsResponse{}, handling.Wrap(errors.New("credentials duplicate found"), handling.WithCode(codes.AlreadyExists))
		}

		domainCredentials, err := ConvertToDomainCredentials(credentials)
		if err != nil {
			log.Error("unable to convert to domain credentials", logging.Error(err))
			return business.CreateCredentialsResponse{}, handling.Wrap(errors.New("error generating password hash"), handling.WithCode(codes.Internal))
		}

		if err := service.repository.Save(ctx, domainCredentials); err != nil {
			log.Error("error while saving credentials to a database", logging.Error(err))
			return business.CreateCredentialsResponse{}, handling.Process(err)
		}

		log.Info("credentials successfully created")
		return business.CreateCredentialsResponse{
			CreatedCredentialsID: domainCredentials.ID,
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("credentials creating deadline exceeded")
			return business.CreateCredentialsResponse{}, handling.Wrap(errCreateCredentialsDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while creating credentials", logging.Error(err))
		return business.CreateCredentialsResponse{}, handling.Process(err)
	}

	log.Info("credentials successfully created", slog.Any("createdCredentialsID", credentialsCreateResponse.CreatedCredentialsID))
	return credentialsCreateResponse, nil
}
