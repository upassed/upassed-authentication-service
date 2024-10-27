package credentials

import (
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"log/slog"
)

const keyFormat = "credentials:%s"

type RedisClient struct {
	cfg    *config.Config
	log    *slog.Logger
	client *redis.Client
}

func New(client *redis.Client, cfg *config.Config, log *slog.Logger) *RedisClient {
	return &RedisClient{
		cfg:    cfg,
		log:    log,
		client: client,
	}
}
