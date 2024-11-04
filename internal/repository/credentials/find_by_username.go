package credentials

import (
	"context"
	"errors"
	"github.com/upassed/upassed-authentication-service/internal/handling"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	"github.com/upassed/upassed-authentication-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

var (
	errSearchingCredentialsByUsername = errors.New("error while searching credentials by id")
	ErrCredentialsNotFoundByUsername  = errors.New("credentials by username not found in database")
)

func (repository *credentialsRepositoryImpl) FindByUsername(ctx context.Context, username string) (*domain.Credentials, error) {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.CredentialsTracerName).Start(ctx, "credentialsRepository#FindByUsername")
	span.SetAttributes(attribute.String("username", username))
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.FindByUsername),
		logging.WithCtx(ctx),
		logging.WithAny("username", username),
	)

	log.Info("started searching credentials by username in redis cache")
	credentialsFromCache, err := repository.cache.GetByUsername(spanContext, username)
	if err == nil {
		log.Info("credentials was found in cache, not going to the database")
		return credentialsFromCache, nil
	}

	log.Info("started searching credentials in a database")
	foundCredentials := domain.Credentials{}
	searchResult := repository.db.WithContext(ctx).Where("username = ?", username).First(&foundCredentials)
	if err := searchResult.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("credentials was not found in the database", logging.Error(err))
			tracing.SetSpanError(span, err)
			return nil, handling.New(ErrCredentialsNotFoundByUsername.Error(), codes.NotFound)
		}

		log.Error("error while searching credentials in the database", logging.Error(err))
		tracing.SetSpanError(span, err)
		return nil, handling.New(errSearchingCredentialsByUsername.Error(), codes.Internal)
	}

	log.Info("credentials were successfully found in a database")
	log.Info("saving credentials to cache")
	if err := repository.cache.Save(spanContext, &foundCredentials); err != nil {
		log.Error("error while saving credentials to cache", logging.Error(err))
		tracing.SetSpanError(span, err)
	}

	log.Info("credentials were successfully saved to the cache")
	return &foundCredentials, nil
}
