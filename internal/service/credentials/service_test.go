package credentials_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/service/credentials"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"github.com/upassed/upassed-authentication-service/internal/util/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var (
	cfg        *config.Config
	repository *mocks.CredentialsRepository
	service    credentials.Service
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

	repository = mocks.NewCredentialsRepository(ctrl)
	service = credentials.New(cfg, logging.New(config.EnvTesting), repository)

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCreate_ErrorCheckingDuplicateExists(t *testing.T) {
	credentialsToCreate := util.RandomBusinessCredentials()
	expectedRepositoryError := errors.New("some repo error")

	repository.EXPECT().
		CheckDuplicatesExists(gomock.Any(), credentialsToCreate.Username).
		Return(false, expectedRepositoryError)

	_, err := service.Create(context.Background(), credentialsToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
	assert.Equal(t, codes.Internal, convertedError.Code())
}

func TestCreate_DuplicateExists(t *testing.T) {
	credentialsToCreate := util.RandomBusinessCredentials()
	repository.EXPECT().
		CheckDuplicatesExists(gomock.Any(), credentialsToCreate.Username).
		Return(true, nil)

	_, err := service.Create(context.Background(), credentialsToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, "credentials duplicate found", convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())
}

func TestCreate_ErrorSavingCredentialsToDatabase(t *testing.T) {
	credentialsToCreate := util.RandomBusinessCredentials()
	repository.EXPECT().
		CheckDuplicatesExists(gomock.Any(), credentialsToCreate.Username).
		Return(false, nil)

	expectedRepositoryError := errors.New("err while saving credentials")
	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(expectedRepositoryError)

	_, err := service.Create(context.Background(), credentialsToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedRepositoryError.Error(), convertedError.Message())
}

func TestCreate_DeadlineExceeded(t *testing.T) {
	oldTimeout := cfg.Timeouts.EndpointExecutionTimeoutMS
	cfg.Timeouts.EndpointExecutionTimeoutMS = "0"

	credentialsToCreate := util.RandomBusinessCredentials()
	repository.EXPECT().
		CheckDuplicatesExists(gomock.Any(), credentialsToCreate.Username).
		Return(false, nil)

	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	_, err := service.Create(context.Background(), credentialsToCreate)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, credentials.ErrCreateCredentialsDeadlineExceeded.Error(), convertedError.Message())

	cfg.Timeouts.EndpointExecutionTimeoutMS = oldTimeout
}

func TestCreate_HappyPath(t *testing.T) {
	credentialsToCreate := util.RandomBusinessCredentials()
	repository.EXPECT().
		CheckDuplicatesExists(gomock.Any(), credentialsToCreate.Username).
		Return(false, nil)

	repository.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Return(nil)

	response, err := service.Create(context.Background(), credentialsToCreate)
	require.NoError(t, err)

	assert.Equal(t, credentialsToCreate.ID, response.CreatedCredentialsID)
}
