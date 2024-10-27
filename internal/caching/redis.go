package caching

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"log/slog"
	"net"
	"strconv"
)

var (
	errCreatingRedisClient = errors.New("unable to create a redis client")
)

func OpenRedisConnection(cfg *config.Config, log *slog.Logger) (*redis.Client, error) {
	log = logging.Wrap(log,
		logging.WithOp(OpenRedisConnection),
	)

	log.Info("started creating redis connection")
	databaseNumber, err := strconv.Atoi(cfg.Redis.DatabaseNumber)
	if err != nil {
		log.Error("unable to parse redis database number", logging.Error(err))
		return nil, err
	}

	redisDatabase := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port),
		Username: cfg.Redis.User,
		Password: cfg.Redis.DatabaseNumber,
		DB:       databaseNumber,
	})

	log.Info("pinging redis database")
	if _, err := redisDatabase.Ping(context.Background()).Result(); err != nil {
		log.Error("unable to ping redis database", logging.Error(err))
		return nil, errCreatingRedisClient
	}

	log.Info("redis client successfully created")
	return redisDatabase, nil
}
