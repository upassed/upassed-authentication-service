package credentials

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

var (
	ErrCredentialsIsNotPresentInCache     = errors.New("credentials id is not present as a key in the cache")
	errFetchingCredentialsFromCache       = errors.New("unable to get credentials by id from the cache")
	errUnmarshallingCredentialsDataToJson = errors.New("unable to unmarshall credentials data from the cache to json format")
)

func (client *RedisClient) GetByID(ctx context.Context, credentialsID uuid.UUID) (*domain.Credentials, error) {
	_, span := otel.Tracer(client.cfg.Tracing.CredentialsTracerName).Start(ctx, "redisClient#GetByID")
	span.SetAttributes(attribute.String("id", credentialsID.String()))
	defer span.End()

	log := logging.Wrap(client.log,
		logging.WithOp(client.GetByID),
		logging.WithCtx(ctx),
		logging.WithAny("credentialsID", credentialsID),
	)

	log.Info("started getting credentials data by id from cache")
	credentialsData, err := client.client.Get(ctx, fmt.Sprintf(keyFormat, credentialsID.String())).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.Error("credentials by id was not found in cache")
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, ErrCredentialsIsNotPresentInCache
		}

		log.Error("error while fetching credentials by id from cache", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, errFetchingCredentialsFromCache
	}

	log.Info("credentials by id was found in cache, unmarshalling from json")
	var credentials domain.Credentials
	if err := json.Unmarshal([]byte(credentialsData), &credentials); err != nil {
		log.Error("error while unmarshalling credentials data to json", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, errUnmarshallingCredentialsDataToJson
	}

	log.Info("credentials was successfully found and unmarshalled")
	return &credentials, nil
}
