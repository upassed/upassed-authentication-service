package credentials

import (
	"context"
	"github.com/upassed/upassed-authentication-service/internal/config"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	"gorm.io/gorm"
	"log/slog"
)

type Repository interface {
	CheckDuplicatesExists(ctx context.Context, username string) (bool, error)
	Save(context.Context, *domain.Credentials) error
}

type credentialsRepositoryImpl struct {
	db  *gorm.DB
	cfg *config.Config
	log *slog.Logger
}

func New(db *gorm.DB, cfg *config.Config, log *slog.Logger) Repository {
	return &credentialsRepositoryImpl{
		db:  db,
		cfg: cfg,
		log: log,
	}
}
