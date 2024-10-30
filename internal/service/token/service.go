package token

import (
	"context"
	"errors"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/jwt"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"log/slog"
)

var (
	ErrParsingToken            = errors.New("unable to parse token")
	ErrGeneratingTokens        = errors.New("error while generating access and refresh tokens")
	ErrTokenInvalid            = errors.New("token is invalid or expired")
	errExtractingTokenClaims   = errors.New("unable to extract map claims from token")
	errUsernameClaimNotPresent = errors.New("username key is not present in refresh token claims")
)

type Service interface {
	Generate(context.Context, *business.TokenGenerateRequest) (*business.TokenGenerateResponse, error)
	Refresh(context.Context, *business.TokenRefreshRequest) (*business.TokenRefreshResponse, error)
	Validate(context.Context, *business.TokenValidateRequest) (*business.TokenValidateResponse, error)
}

type tokenServiceImpl struct {
	cfg                   *config.Config
	log                   *slog.Logger
	tokenGenerator        jwt.TokenGenerator
	credentialsRepository credentialsRepository
}

type credentialsRepository interface {
	FindByUsername(ctx context.Context, username string) (*domain.Credentials, error)
}

func New(cfg *config.Config, log *slog.Logger, tokenGenerator jwt.TokenGenerator, credentialsRepository credentialsRepository) Service {
	return &tokenServiceImpl{
		cfg:                   cfg,
		log:                   log,
		tokenGenerator:        tokenGenerator,
		credentialsRepository: credentialsRepository,
	}
}
