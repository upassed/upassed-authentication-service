package token

import (
	"context"
	"github.com/upassed/upassed-authentication-service/internal/config"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"log/slog"
)

type Service interface {
	Generate(context.Context, *business.TokenGenerateRequest) (*business.TokenGenerateResponse, error)
	Refresh(context.Context, *business.TokenRefreshRequest) (*business.TokenRefreshResponse, error)
	Validate(context.Context, *business.TokenValidateRequest) (*business.TokenValidateResponse, error)
}

type tokenServiceImpl struct {
	cfg                   *config.Config
	log                   *slog.Logger
	credentialsRepository credentialsRepository
}

type credentialsRepository interface {
	FindByUsername(ctx context.Context, username string) (*domain.Credentials, error)
}

func New(cfg *config.Config, log *slog.Logger, credentialsRepository credentialsRepository) Service {
	return &tokenServiceImpl{
		cfg:                   cfg,
		log:                   log,
		credentialsRepository: credentialsRepository,
	}
}
