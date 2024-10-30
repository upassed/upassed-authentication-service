package token_test

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/jwt"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
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
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
}

func TestCreate_FindingCredentialsInRepoDeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	request := util.RandomBusinessTokenGenerateRequest()
	foundCredentials := util.RandomDomainCredentials()

	credentialsRepository := new(mockCredentialsRepository)
	credentialsRepository.On("FindByUsername", mock.Anything, request.Username).Return(foundCredentials, nil)

	logger := logging.New(config.EnvTesting)
	service := token.New(cfg, logger, new(mockTokenGenerator), credentialsRepository)
	_, err := service.Generate(context.Background(), request)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, token.ErrFindingCredentialsByUsernameDeadlineExceeded.Error(), convertedError.Message())
	assert.Equal(t, codes.DeadlineExceeded, convertedError.Code())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestCreate_PasswordHashNotMatch(t *testing.T) {
	request := util.RandomBusinessTokenGenerateRequest()
	foundCredentials := util.RandomDomainCredentials()

	credentialsRepository := new(mockCredentialsRepository)
	credentialsRepository.On("FindByUsername", mock.Anything, request.Username).Return(foundCredentials, nil)

	logger := logging.New(config.EnvTesting)
	service := token.New(cfg, logger, new(mockTokenGenerator), credentialsRepository)
	_, err := service.Generate(context.Background(), request)
	require.NotNil(t, err)

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
	require.NotNil(t, err)

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
