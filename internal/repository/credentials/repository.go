package credentials

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/upassed/upassed-authentication-service/internal/caching/credentials"
	"github.com/upassed/upassed-authentication-service/internal/config"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	"gorm.io/gorm"
	"log/slog"
)

type Repository interface {
	CheckDuplicatesExists(ctx context.Context, username string) (bool, error)
	Save(context.Context, *domain.Credentials) error
	FindByUsername(ctx context.Context, username string) (*domain.Credentials, error)
}

type repositoryImpl struct {
	db    *gorm.DB
	cache *credentials.RedisClient
	cfg   *config.Config
	log   *slog.Logger
}

func New(db *gorm.DB, redisClient *redis.Client, cfg *config.Config, log *slog.Logger) Repository {
	cacheClient := credentials.New(redisClient, cfg, log)

	return &repositoryImpl{
		db:    db,
		cache: cacheClient,
		cfg:   cfg,
		log:   log,
	}
}
