package credentials_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	"github.com/upassed/upassed-authentication-service/internal/service/credentials"
	"github.com/upassed/upassed-authentication-service/internal/util"
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

func (m *mockCredentialsRepository) CheckDuplicatesExists(ctx context.Context, username string) (bool, error) {
	args := m.Called(ctx, username)
	return args.Bool(0), args.Error(1)
}

func (m *mockCredentialsRepository) Save(ctx context.Context, credentials *domain.Credentials) error {
	args := m.Called(ctx, credentials)
	return args.Error(0)
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

func TestCreate_ErrorCheckingDuplicateExists(t *testing.T) {
	credentialsRepository := new(mockCredentialsRepository)

	credentialsToCreate := util.RandomBusinessCredentials()
	expectedRepositoryError := errors.New("some repo error")
	credentialsRepository.On(
		"CheckDuplicatesExists",
		mock.Anything,
		credentialsToCreate.Username,
	).Return(false, expectedRepositoryError)

	service := credentials.New(cfg, logging.New(config.EnvTesting), credentialsRepository)
	_, err := service.Create(context.Background(), credentialsToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_DuplicateExists(t *testing.T) {
	credentialsRepository := new(mockCredentialsRepository)

	credentialsToCreate := util.RandomBusinessCredentials()
	credentialsRepository.On(
		"CheckDuplicatesExists",
		mock.Anything,
		credentialsToCreate.Username,
	).Return(true, nil)

	service := credentials.New(cfg, logging.New(config.EnvTesting), credentialsRepository)
	_, err := service.Create(context.Background(), credentialsToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "credentials duplicate found", convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TestCreate_ErrorSavingCredentialsToDatabase(t *testing.T) {
	credentialsRepository := new(mockCredentialsRepository)

	credentialsToCreate := util.RandomBusinessCredentials()
	credentialsRepository.On(
		"CheckDuplicatesExists",
		mock.Anything,
		credentialsToCreate.Username,
	).Return(false, nil)

	expectedRepositoryError := errors.New("err while saving credentials")
	credentialsRepository.On(
		"Save",
		mock.Anything,
		mock.Anything,
	).Return(expectedRepositoryError)

	service := credentials.New(cfg, logging.New(config.EnvTesting), credentialsRepository)
	_, err := service.Create(context.Background(), credentialsToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
}

func TestCreate_DeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	credentialsRepository := new(mockCredentialsRepository)

	credentialsToCreate := util.RandomBusinessCredentials()
	credentialsRepository.On(
		"CheckDuplicatesExists",
		mock.Anything,
		credentialsToCreate.Username,
	).Return(false, nil)

	credentialsRepository.On(
		"Save",
		mock.Anything,
		mock.Anything,
	).Return(nil)

	service := credentials.New(cfg, logging.New(config.EnvTesting), credentialsRepository)
	_, err := service.Create(context.Background(), credentialsToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, credentials.ErrCreateCredentialsDeadlineExceeded.Error(), convertedError.Message())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestCreate_HappyPath(t *testing.T) {
	credentialsRepository := new(mockCredentialsRepository)

	credentialsToCreate := util.RandomBusinessCredentials()
	credentialsRepository.On(
		"CheckDuplicatesExists",
		mock.Anything,
		credentialsToCreate.Username,
	).Return(false, nil)

	credentialsRepository.On(
		"Save",
		mock.Anything,
		mock.Anything,
	).Return(nil)

	service := credentials.New(cfg, logging.New(config.EnvTesting), credentialsRepository)
	response, err := service.Create(context.Background(), credentialsToCreate)
	require.NoError(t, err)

	assert.Equal(t, credentialsToCreate.ID, response.CreatedCredentialsID)
}
