package token_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/server"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
	"testing"
)

type mockTokenService struct {
	mock.Mock
}

func (m *mockTokenService) Generate(ctx context.Context, request *business.TokenGenerateRequest) (*business.TokenGenerateResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*business.TokenGenerateResponse), args.Error(1)
}

func (m *mockTokenService) Refresh(ctx context.Context, request *business.TokenRefreshRequest) (*business.TokenRefreshResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*business.TokenRefreshResponse), args.Error(1)
}

func (m *mockTokenService) Validate(ctx context.Context, request *business.TokenValidateRequest) (*business.TokenValidateResponse, error) {
	args := m.Called(ctx, request)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*business.TokenValidateResponse), args.Error(1)
}

var (
	tokenClient client.TokenClient
	tokenSvc    *mockTokenService
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

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("cfg load error: ", err)
	}

	logger := logging.New(cfg.Env)
	tokenSvc = new(mockTokenService)
	tokenServer := server.New(server.AppServerCreateParams{
		Config:       cfg,
		Log:          logger,
		TokenService: tokenSvc,
	})

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	cc, err := grpc.NewClient(fmt.Sprintf(":%s", cfg.GrpcServer.Port), opts...)
	if err != nil {
		log.Fatal("error creating client connection", err)
	}

	tokenClient = client.NewTokenClient(cc)
	go func() {
		if err := tokenServer.Run(); err != nil {
			os.Exit(1)
		}
	}()

	exitCode := m.Run()
	tokenServer.GracefulStop()
	os.Exit(exitCode)
}

func TestGenerateToken_InvalidRequest(t *testing.T) {
	request := util.RandomClientTokenGenerateRequest()
	request.Username = "_invalid_"

	_, err := tokenClient.Generate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.InvalidArgument, convertedError.Code())

	clearTokenServiceMockCalls()
}

func TestGenerateToken_ServiceError(t *testing.T) {
	request := util.RandomClientTokenGenerateRequest()

	expectedServiceError := errors.New("some service error")
	tokenSvc.On("Generate", mock.Anything, mock.Anything).Return(nil, expectedServiceError)

	_, err := tokenClient.Generate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedServiceError.Error(), convertedError.Message())

	clearTokenServiceMockCalls()
}

func TestGenerateToken_HappyPath(t *testing.T) {
	request := util.RandomClientTokenGenerateRequest()

	expectedServiceResponse := util.RandomBusinessTokenGenerateResponse()
	tokenSvc.On("Generate", mock.Anything, mock.Anything).Return(expectedServiceResponse, nil)

	response, err := tokenClient.Generate(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, expectedServiceResponse.AccessToken, response.GetAccessToken())
	assert.Equal(t, expectedServiceResponse.RefreshToken, response.GetRefreshToken())

	clearTokenServiceMockCalls()
}

func TestRefreshToken_ServiceError(t *testing.T) {
	request := util.RandomClientTokenRefreshRequest()

	expectedServiceError := errors.New("some service error")
	tokenSvc.On("Refresh", mock.Anything, mock.Anything).Return(nil, expectedServiceError)

	_, err := tokenClient.Refresh(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedServiceError.Error(), convertedError.Message())

	clearTokenServiceMockCalls()
}

func TestRefreshToken_HappyPath(t *testing.T) {
	request := util.RandomClientTokenRefreshRequest()

	expectedServiceResponse := util.RandomBusinessTokenRefreshResponse()
	tokenSvc.On("Refresh", mock.Anything, mock.Anything).Return(expectedServiceResponse, nil)

	response, err := tokenClient.Refresh(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, expectedServiceResponse.NewAccessToken, response.GetNewAccessToken())

	clearTokenServiceMockCalls()
}

func TestValidateToken_ServiceError(t *testing.T) {
	request := util.RandomClientTokenValidateRequest()

	expectedServiceError := errors.New("some service error")
	tokenSvc.On("Validate", mock.Anything, mock.Anything).Return(nil, expectedServiceError)

	_, err := tokenClient.Validate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedServiceError.Error(), convertedError.Message())

	clearTokenServiceMockCalls()
}

func TestValidateToken_HappyPath(t *testing.T) {
	request := util.RandomClientTokenValidateRequest()

	expectedServiceResponse := util.RandomBusinessTokenValidateResponse()
	tokenSvc.On("Validate", mock.Anything, mock.Anything).Return(expectedServiceResponse, nil)

	response, err := tokenClient.Validate(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, expectedServiceResponse.Username, response.Username)
	assert.Equal(t, string(expectedServiceResponse.AccountType), response.GetAccountType())

	clearTokenServiceMockCalls()
}

func clearTokenServiceMockCalls() {
	tokenSvc.Calls = nil
	tokenSvc.ExpectedCalls = nil
}
