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
	"log/slog"
)

var (
	errCountingDuplicatesCredentials = errors.New("error while counting duplicate credentials")
)

func (repository *credentialsRepositoryImpl) CheckDuplicatesExists(ctx context.Context, username string) (bool, error) {
	spanContext, span := otel.Tracer(repository.cfg.Tracing.CredentialsTracerName).Start(ctx, "credentialsRepository#CheckDuplicateExists")
	span.SetAttributes(attribute.String("username", username))
	defer span.End()

	log := logging.Wrap(repository.log,
		logging.WithOp(repository.CheckDuplicatesExists),
		logging.WithCtx(ctx),
		logging.WithAny("username", username),
	)

	log.Info("started searching credentials by username in redis cache")
	_, err := repository.cache.GetByUsername(spanContext, username)
	if err == nil {
		log.Info("credentials was found in cache, not going to the database")
		return true, nil
	}

	log.Info("started checking credentials duplicates by username")
	var credentialsCount int64
	countResult := repository.db.WithContext(ctx).Model(&domain.Credentials{}).Where("username = ?", username).Count(&credentialsCount)
	if err := countResult.Error; err != nil {
		log.Error("error while counting credentials with username in database")
		span.SetAttributes(attribute.String("err", err.Error()))
		return false, handling.New(errCountingDuplicatesCredentials.Error(), codes.Internal)
	}

	if credentialsCount > 0 {
		log.Info("found credentials duplicates in database", slog.Int64("duplicatesCount", credentialsCount))
		return true, nil
	}

	log.Info("credentials duplicates not found in database")
	return false, nil
}
