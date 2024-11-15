package token_test

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/jwt"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/internal/service/token"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"github.com/upassed/upassed-authentication-service/internal/util/mocks"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	cfg            *config.Config
	tokenGenerator *mocks.TokenGenerator
	repository     *mocks.CredentialsRepository
	service        token.Service
)

func TestMain(m *testing.M) {
	currentDir, _ := os.Getwd()
	projectRoot, err := util.GetProjectRoot(currentDir)
	if err != nil {
		log.Fatal("error to get project root folder: ", err)
	}

	if err := os.Setenv(config.EnvConfigPath, filepath.Join(projectRoot, "config", "test.yml")); err != nil {
		log.Fatal(err)
	}

	cfg, err = config.Load()
	if err != nil {
		log.Fatal("unable to parse config: ", err)
	}

	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	tokenGenerator = mocks.NewTokenGenerator(ctrl)
	repository = mocks.NewCredentialsRepository(ctrl)
	service = token.New(cfg, logging.New(config.EnvTesting), tokenGenerator, repository)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreate_ErrorFindingCredentialsByUsername(t *testing.T) {
	request := util.RandomBusinessTokenGenerateRequest()

	expectedRepositoryError := errors.New("some repo error")
	repository.EXPECT().
		FindByUsername(gomock.Any(), request.Username).
		Return(nil, expectedRepositoryError)

	_, err := service.Generate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
}

func TestCreate_PasswordHashNotMatch(t *testing.T) {
	request := util.RandomBusinessTokenGenerateRequest()
	foundCredentials := util.RandomDomainCredentials()

	repository.EXPECT().
		FindByUsername(gomock.Any(), request.Username).
		Return(foundCredentials, nil)

	_, err := service.Generate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, token.ErrPasswordHashNotMatch.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_ErrorGeneratingTokens(t *testing.T) {
	request := util.RandomBusinessTokenGenerateRequest()
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	foundCredentials := &domain.Credentials{
		ID:           uuid.New(),
		Username:     request.Username,
		PasswordHash: hash,
	}

	repository.EXPECT().
		FindByUsername(gomock.Any(), request.Username).
		Return(foundCredentials, nil)

	tokenGenerator.EXPECT().
		GenerateFor(request.Username).
		Return(nil, errors.New("some token generator error"))

	_, err = service.Generate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, token.ErrGeneratingTokens.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_HappyPath(t *testing.T) {
	request := util.RandomBusinessTokenGenerateRequest()
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	require.NoError(t, err)

	foundCredentials := &domain.Credentials{
		ID:           uuid.New(),
		Username:     request.Username,
		PasswordHash: hash,
	}

	repository.EXPECT().
		FindByUsername(gomock.Any(), request.Username).
		Return(foundCredentials, nil)

	generatedTokens := util.RandomJwtGeneratedTokens()
	tokenGenerator.EXPECT().
		GenerateFor(request.Username).
		Return(generatedTokens, nil)

	response, err := service.Generate(context.Background(), request)
	require.NoError(t, err)

	assert.NotNil(t, response.AccessToken)
	assert.NotNil(t, response.RefreshToken)
}

func TestRefresh_InvalidRefreshToken(t *testing.T) {
	request := util.RandomBusinessTokenRefreshRequest()

	_, err := service.Refresh(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, token.ErrParsingToken.Error(), convertedError.Message())
}

func TestRefresh_ExpiredRefreshToken(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	oldJwtRefreshTokenTTL := cfg.Jwt.RefreshTokenTTL
	cfg.Jwt.RefreshTokenTTL = "-10m"

	username := gofakeit.Username()
	tokens, err := jwt.New(cfg, logger).GenerateFor(username)
	require.NoError(t, err)

	request := &business.TokenRefreshRequest{
		RefreshToken: tokens.RefreshToken,
	}

	_, err = service.Refresh(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, token.ErrParsingToken.Error(), convertedError.Message())

	cfg.Jwt.RefreshTokenTTL = oldJwtRefreshTokenTTL
}

func TestRefresh_TokenGenerationError(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	username := gofakeit.Username()
	generator := jwt.New(cfg, logger)
	tokens, err := generator.GenerateFor(username)
	require.NoError(t, err)

	request := &business.TokenRefreshRequest{
		RefreshToken: tokens.RefreshToken,
	}

	tokenGenerator.EXPECT().
		GenerateFor(username).
		Return(nil, errors.New("some error"))

	_, err = service.Refresh(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, token.ErrGeneratingTokens.Error(), convertedError.Message())
}

func TestRefresh_HappyPath(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	username := gofakeit.Username()
	generator := jwt.New(cfg, logger)
	tokens, err := generator.GenerateFor(username)
	require.NoError(t, err)

	request := &business.TokenRefreshRequest{
		RefreshToken: tokens.RefreshToken,
	}

	service := token.New(cfg, logger, generator, repository)
	response, err := service.Refresh(context.Background(), request)
	require.NoError(t, err)

	assert.NotNil(t, response.NewAccessToken)
}

func TestValidate_InvalidAccessToken(t *testing.T) {
	request := util.RandomBusinessTokenValidateRequest()

	_, err := service.Validate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, token.ErrParsingToken.Error(), convertedError.Message())
}

func TestValidate_ExpiredAccessToken(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	oldJwtAccessTokenTTL := cfg.Jwt.AccessTokenTTL
	cfg.Jwt.AccessTokenTTL = "-10m"

	username := gofakeit.Username()
	tokens, err := jwt.New(cfg, logger).GenerateFor(username)
	require.NoError(t, err)

	request := &business.TokenValidateRequest{
		AccessToken: tokens.AccessToken,
	}

	_, err = service.Validate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, token.ErrParsingToken.Error(), convertedError.Message())

	cfg.Jwt.AccessTokenTTL = oldJwtAccessTokenTTL
}

func TestValidate_HappyPath(t *testing.T) {
	logger := logging.New(config.EnvTesting)
	username := gofakeit.Username()
	generator := jwt.New(cfg, logger)
	tokens, err := generator.GenerateFor(username)
	require.NoError(t, err)

	request := &business.TokenValidateRequest{
		AccessToken: tokens.AccessToken,
	}

	foundCredentials := util.RandomDomainCredentials()
	foundCredentials.Username = username

	repository.EXPECT().
		FindByUsername(gomock.Any(), username).
		Return(foundCredentials, nil)

	response, err := service.Validate(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, foundCredentials.ID, response.CredentialsID)
	assert.Equal(t, username, response.Username)
	assert.Equal(t, business.AccountType(foundCredentials.AccountType), response.AccountType)
}
