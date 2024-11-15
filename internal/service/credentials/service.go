package credentials

import (
	"context"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"log/slog"

	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
)

type Service interface {
	Create(context.Context, *business.Credentials) (*business.CreateCredentialsResponse, error)
}

type serviceImpl struct {
	cfg        *config.Config
	log        *slog.Logger
	repository repository
}

type repository interface {
	CheckDuplicatesExists(ctx context.Context, username string) (bool, error)
	Save(context.Context, *domain.Credentials) error
}

func New(cfg *config.Config, log *slog.Logger, repository repository) Service {
	return &serviceImpl{
		cfg:        cfg,
		log:        log,
		repository: repository,
	}
}
