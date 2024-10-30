package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"log/slog"
	"time"
)

var (
	errGeneratingAccessToken  = errors.New("error while generating access token")
	errGeneratingRefreshToken = errors.New("error while generating refresh token")
)

type GeneratedTokens struct {
	AccessToken  string
	RefreshToken string
}

type GenerateTokensParams struct {
	Username        string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Secret          string
}

type TokenGenerator interface {
	GenerateFor(username string) (*GeneratedTokens, error)
}

type tokenGeneratorImpl struct {
	cfg *config.Config
	log *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) TokenGenerator {
	return &tokenGeneratorImpl{
		cfg: cfg,
		log: log,
	}
}

func (generator *tokenGeneratorImpl) GenerateFor(username string) (*GeneratedTokens, error) {
	accessToken, err := generator.generateAccessToken(username)
	if err != nil {
		return nil, errGeneratingAccessToken
	}

	refreshToken, err := generator.generateRefreshToken(username)
	if err != nil {
		return nil, errGeneratingRefreshToken
	}

	return &GeneratedTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

const (
	UsernameKey = "username"
	ExpKey      = "exp"
)

func (generator *tokenGeneratorImpl) generateAccessToken(username string) (string, error) {
	claims := jwt.MapClaims{
		UsernameKey: username,
		ExpKey:      time.Now().Add(generator.cfg.GetJwtAccessTokenTTL()).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(generator.cfg.Jwt.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (generator *tokenGeneratorImpl) generateRefreshToken(username string) (string, error) {
	claims := jwt.MapClaims{
		UsernameKey: username,
		ExpKey:      time.Now().Add(generator.cfg.GetJwtRefreshTokenTTL()).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(generator.cfg.Jwt.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
