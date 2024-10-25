package credentials

import (
	"context"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"log/slog"

	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
)

type Service interface {
	CreateCredentials(context.Context, business.Credentials) (business.CreateCredentialsResponse, error)
}

type credentialsServiceImpl struct {
	cfg        *config.Config
	log        *slog.Logger
	repository authenticationRepository
}

type authenticationRepository interface {
	SaveCredentials(context.Context, domain.Credentials) error
}

func New(cfg *config.Config, log *slog.Logger, repository authenticationRepository) Service {
	return &credentialsServiceImpl{
		cfg:        cfg,
		log:        log,
		repository: repository,
	}
}
