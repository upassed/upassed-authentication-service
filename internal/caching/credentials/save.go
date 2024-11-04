package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	"github.com/upassed/upassed-authentication-service/internal/tracing"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	errMarshallingCredentialsData   = errors.New("unable to marshall credentials data to json format")
	errSavingCredentialsDataToCache = errors.New("unable to save credentials data to redis cache")
)

func (client *RedisClient) Save(ctx context.Context, credentials *domain.Credentials) error {
	_, span := otel.Tracer(client.cfg.Tracing.CredentialsTracerName).Start(ctx, "redisClient#Save")
	span.SetAttributes(attribute.String("username", credentials.Username))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.Save),
		logging.WithCtx(ctx),
		logging.WithAny("username", credentials.Username),
	)

	log.Info("marshalling credentials data to json to save to cache")
	jsonCredentialsData, err := json.Marshal(credentials)
	if err != nil {
		log.Error("unable to marshall credentials data to json format")
		tracing.SetSpanError(span, err)
		return errMarshallingCredentialsData
	}

	log.Info("saving credentials data to the cache")
	if err := client.client.Set(ctx, fmt.Sprintf(keyFormat, credentials.Username), jsonCredentialsData, client.cfg.GetRedisEntityTTL()).Err(); err != nil {
		log.Error("error while saving credentials data to the cache", logging.Error(err))
		tracing.SetSpanError(span, err)
		return errSavingCredentialsDataToCache
	}

	log.Info("credentials successfully saved to the cache")
	return nil
}
