package token_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/server"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"github.com/upassed/upassed-authentication-service/internal/util/mocks"
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

var (
	tokenClient client.TokenClient
	tokenSvc    *mocks.TokenService
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

	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	tokenSvc = mocks.NewTokenService(ctrl)
	tokenServer := server.New(server.AppServerCreateParams{
		Config:       cfg,
		Log:          logging.New(config.EnvTesting),
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
}

func TestGenerateToken_ServiceError(t *testing.T) {
	request := util.RandomClientTokenGenerateRequest()

	expectedServiceError := errors.New("some service error")
	tokenSvc.EXPECT().
		Generate(gomock.Any(), gomock.Any()).
		Return(nil, expectedServiceError)

	_, err := tokenClient.Generate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedServiceError.Error(), convertedError.Message())
}

func TestGenerateToken_HappyPath(t *testing.T) {
	request := util.RandomClientTokenGenerateRequest()

	expectedServiceResponse := util.RandomBusinessTokenGenerateResponse()
	tokenSvc.EXPECT().
		Generate(gomock.Any(), gomock.Any()).
		Return(expectedServiceResponse, nil)

	response, err := tokenClient.Generate(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, expectedServiceResponse.AccessToken, response.GetAccessToken())
	assert.Equal(t, expectedServiceResponse.RefreshToken, response.GetRefreshToken())
}

func TestRefreshToken_ServiceError(t *testing.T) {
	request := util.RandomClientTokenRefreshRequest()

	expectedServiceError := errors.New("some service error")
	tokenSvc.EXPECT().
		Refresh(gomock.Any(), gomock.Any()).
		Return(nil, expectedServiceError)

	_, err := tokenClient.Refresh(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedServiceError.Error(), convertedError.Message())
}

func TestRefreshToken_HappyPath(t *testing.T) {
	request := util.RandomClientTokenRefreshRequest()

	expectedServiceResponse := util.RandomBusinessTokenRefreshResponse()
	tokenSvc.EXPECT().
		Refresh(gomock.Any(), gomock.Any()).
		Return(expectedServiceResponse, nil)

	response, err := tokenClient.Refresh(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, expectedServiceResponse.NewAccessToken, response.GetNewAccessToken())
}

func TestValidateToken_ServiceError(t *testing.T) {
	request := util.RandomClientTokenValidateRequest()

	expectedServiceError := errors.New("some service error")
	tokenSvc.EXPECT().
		Validate(gomock.Any(), gomock.Any()).
		Return(nil, expectedServiceError)

	_, err := tokenClient.Validate(context.Background(), request)
	require.Error(t, err)

	convertedError := status.Convert(err)
	assert.Equal(t, expectedServiceError.Error(), convertedError.Message())
}

func TestValidateToken_HappyPath(t *testing.T) {
	request := util.RandomClientTokenValidateRequest()

	expectedServiceResponse := util.RandomBusinessTokenValidateResponse()
	tokenSvc.EXPECT().
		Validate(gomock.Any(), gomock.Any()).
		Return(expectedServiceResponse, nil)

	response, err := tokenClient.Validate(context.Background(), request)
	require.NoError(t, err)

	assert.Equal(t, expectedServiceResponse.Username, response.Username)
	assert.Equal(t, string(expectedServiceResponse.AccountType), response.GetAccountType())
}
