package credentials

import (
	"context"
	"errors"
	"github.com/upassed/upassed-authentication-service/internal/handling"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

var (
	ErrSavingCredentials = errors.New("error while saving credentials")
)

func (repository *credentialsRepositoryImpl) Save(ctx context.Context, credentials *domain.Credentials) error {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.CredentialsTracerName).Start(ctx, "credentialsRepository#Save")
	span.SetAttributes(attribute.String("username", credentials.Username))
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.Save),
		logging.WithCtx(ctx),
		logging.WithAny("credentialsUsername", credentials.Username),
	)

	log.Info("started saving credentials to a database")
	saveResult := repository.db.WithContext(ctx).Create(&credentials)
	if err := saveResult.Error; err != nil || saveResult.RowsAffected != 1 {
		log.Error("error while saving credentials data to a database", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return handling.New(ErrSavingCredentials.Error(), codes.Internal)
	}

	log.Info("credentials were successfully inserted into a database")
	log.Info("saving credentials data into the cache")
	if err := repository.cache.Save(spanContext, credentials); err != nil {
		log.Error("unable to insert credentials in cache", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
	}

	log.Info("credentials were saved to the cache")
	return nil
}
