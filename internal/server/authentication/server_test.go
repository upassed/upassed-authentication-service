package authentication_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/handling"
	"github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/server"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type mockAuthenticationService struct {
	mock.Mock
}

func (m *mockAuthenticationService) CreateCredentials(ctx context.Context, credentials business.Credentials) (business.CreateCredentialsResponse, error) {
	args := m.Called(ctx, credentials)
	return args.Get(0).(business.CreateCredentialsResponse), args.Error(1)
}

var (
	authenticationClient client.AuthenticationClient
	authenticationSvc    *mockAuthenticationService
)

func TestMain(m *testing.M) {
	projectRoot, err := getProjectRoot()
	if err != nil {
		log.Fatal("error to get project root folder: ", err)
	}

	if err := os.Setenv(config.EnvConfigPath, filepath.Join(projectRoot, "config", "test.yml")); err != nil {
		log.Fatal(err)
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal("config load error: ", err)
	}

	logger := logger.New(config.Env)
	authenticationSvc = new(mockAuthenticationService)
	authenticationServer := server.New(server.AppServerCreateParams{
		Config:                config,
		Log:                   logger,
		AuthenticationService: authenticationSvc,
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.NewClient(fmt.Sprintf(":%s", config.GrpcServer.Port), opts...)
	if err != nil {
		log.Fatal("error creating client connection", err)
	}

	authenticationClient = client.NewAuthenticationClient(cc)
	go func() {
		if err := authenticationServer.Run(); err != nil {
			os.Exit(1)
		}
	}()

	exitCode := m.Run()
	authenticationServer.GracefulStop()
	os.Exit(exitCode)
}

func TestCreateCredentials_ServiceError(t *testing.T) {
	request := client.CredentialsCreateRequest{
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 40),
	}

	expectedError := handling.New("some service error", codes.AlreadyExists)
	authenticationSvc.On("CreateCredentials", mock.Anything, mock.Anything).Return(business.CreateCredentialsResponse{}, handling.Process(expectedError))

	_, err := authenticationClient.CreateCredentials(context.Background(), &request)
	require.NotNil(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedError.Error(), convertedError.Message())
	assert.Equal(t, codes.AlreadyExists, convertedError.Code())

	clearAuthenticationServiceMockCalls()
}

func TestCreateCredentials_HappyPath(t *testing.T) {
	request := client.CredentialsCreateRequest{
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 40),
	}

	expectedResponse := business.CreateCredentialsResponse{
		CredentialsID: uuid.New(),
	}

	authenticationSvc.On("CreateCredentials", mock.Anything, mock.Anything).Return(expectedResponse, nil)

	response, err := authenticationClient.CreateCredentials(context.Background(), &request)
	require.Nil(t, err)

	assert.Equal(t, expectedResponse.CredentialsID.String(), response.CredentialsId)

	clearAuthenticationServiceMockCalls()
}

func clearAuthenticationServiceMockCalls() {
	authenticationSvc.Calls = nil
	authenticationSvc.ExpectedCalls = nil
}

func getProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", errors.New("project root not found")
		}

		dir = parentDir
	}
}
