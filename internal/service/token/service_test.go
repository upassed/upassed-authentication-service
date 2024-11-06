package token_test

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/jwt"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/internal/service/token"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
	"testing"
)

type mockCredentialsRepository struct {
	mock.Mock
}

func (m *mockCredentialsRepository) FindByUsername(ctx context.Context, username string) (*domain.Credentials, error) {
	args := m.Called(ctx, username)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.Credentials), args.Error(1)
}

type mockTokenGenerator struct {
	mock.Mock
}

func (m *mockTokenGenerator) GenerateFor(username string) (*jwt.GeneratedTokens, error) {
	args := m.Called(username)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*jwt.GeneratedTokens), args.Error(1)
}

var (
	cfg *config.Config
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

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreate_ErrorFindingCredentialsByUsername(t *testing.T) {
	request := util.RandomBusinessTokenGenerateRequest()

	credentialsRepository := new(mockCredentialsRepository)
	expectedRepositoryError := errors.New("some repo error")
	credentialsRepository.On("FindByUsername", mock.Anything, request.Username).Return(nil, expectedRepositoryError)

	logger := logging.New(config.EnvTesting)
	service := token.New(cfg, logger, new(mockTokenGenerator), credentialsRepository)
	_, err := service.Generate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
}

func TestCreate_PasswordHashNotMatch(t *testing.T) {
	request := util.RandomBusinessTokenGenerateRequest()
	foundCredentials := util.RandomDomainCredentials()

	credentialsRepository := new(mockCredentialsRepository)
	credentialsRepository.On("FindByUsername", mock.Anything, request.Username).Return(foundCredentials, nil)

	logger := logging.New(config.EnvTesting)
	service := token.New(cfg, logger, new(mockTokenGenerator), credentialsRepository)
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

	credentialsRepository := new(mockCredentialsRepository)
	credentialsRepository.On("FindByUsername", mock.Anything, request.Username).Return(foundCredentials, nil)

	tokenGenerator := new(mockTokenGenerator)
	tokenGenerator.On("GenerateFor", request.Username).Return(nil, errors.New("some token generator error"))

	logger := logging.New(config.EnvTesting)
	service := token.New(cfg, logger, tokenGenerator, credentialsRepository)
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

	credentialsRepository := new(mockCredentialsRepository)
	credentialsRepository.On("FindByUsername", mock.Anything, request.Username).Return(foundCredentials, nil)

	tokenGenerator := new(mockTokenGenerator)
	generatedTokens := util.RandomJwtGeneratedTokens()
	tokenGenerator.On("GenerateFor", request.Username).Return(generatedTokens, nil)

	logger := logging.New(config.EnvTesting)
	service := token.New(cfg, logger, tokenGenerator, credentialsRepository)
	response, err := service.Generate(context.Background(), request)
	require.NoError(t, err)

	assert.NotNil(t, response.AccessToken)
	assert.NotNil(t, response.RefreshToken)
}

func TestRefresh_InvalidRefreshToken(t *testing.T) {
	request := util.RandomBusinessTokenRefreshRequest()

	credentialsRepository := new(mockCredentialsRepository)
	tokenGenerator := new(mockTokenGenerator)

	logger := logging.New(config.EnvTesting)
	service := token.New(cfg, logger, tokenGenerator, credentialsRepository)
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

	credentialsRepository := new(mockCredentialsRepository)
	tokenGenerator := new(mockTokenGenerator)

	service := token.New(cfg, logger, tokenGenerator, credentialsRepository)
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

	credentialsRepository := new(mockCredentialsRepository)
	tokenGenerator := new(mockTokenGenerator)
	tokenGenerator.On("GenerateFor", username).Return(nil, errors.New("some error"))

	service := token.New(cfg, logger, tokenGenerator, credentialsRepository)
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

	credentialsRepository := new(mockCredentialsRepository)

	service := token.New(cfg, logger, generator, credentialsRepository)
	response, err := service.Refresh(context.Background(), request)
	require.NoError(t, err)

	assert.NotNil(t, response.NewAccessToken)
}

func TestValidate_InvalidAccessToken(t *testing.T) {
	request := util.RandomBusinessTokenValidateRequest()

	credentialsRepository := new(mockCredentialsRepository)
	tokenGenerator := new(mockTokenGenerator)

	logger := logging.New(config.EnvTesting)
	service := token.New(cfg, logger, tokenGenerator, credentialsRepository)
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

	credentialsRepository := new(mockCredentialsRepository)
	tokenGenerator := new(mockTokenGenerator)

	service := token.New(cfg, logger, tokenGenerator, credentialsRepository)
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

	credentialsRepository := new(mockCredentialsRepository)
	credentialsRepository.On("FindByUsername", mock.Anything, username).Return(foundCredentials, nil)

	service := token.New(cfg, logger, generator, credentialsRepository)
	response, err := service.Validate(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, username, response.Username)
	assert.Equal(t, business.AccountType(foundCredentials.AccountType), response.AccountType)
}
