package credentials_test

import (
	"context"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/caching"
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/repository"
	"github.com/upassed/upassed-authentication-service/internal/repository/credentials"
	"github.com/upassed/upassed-authentication-service/internal/testcontainer"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var (
	credentialsRepository credentials.Repository
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
		log.Fatal("unable to parse config: ", err)
	}

	ctx := context.Background()
	postgresTestcontainer, err := testcontainer.NewPostgresTestcontainer(ctx)
	if err != nil {
		log.Fatal("unable to create a testcontainer: ", err)
	}

	port, err := postgresTestcontainer.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Storage.Port = strconv.Itoa(port)
	logger := logging.New(cfg.Env)
	if err := postgresTestcontainer.Migrate(cfg, logger); err != nil {
		log.Fatal("unable to run migrations: ", err)
	}

	redisTestcontainer, err := testcontainer.NewRedisTestcontainer(ctx, cfg)
	if err != nil {
		log.Fatal("unable to run redis testcontainer: ", err)
	}

	port, err = redisTestcontainer.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Redis.Port = strconv.Itoa(port)
	db, err := repository.OpenGormDbConnection(cfg, logger)
	if err != nil {
		log.Fatal("unable to open connection to postgres: ", err)
	}

	redis, err := caching.OpenRedisConnection(cfg, logger)
	if err != nil {
		log.Fatal("unable to open connection to redis: ", err)
	}

	credentialsRepository = credentials.New(db, redis, cfg, logger)
	exitCode := m.Run()
	if err := postgresTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop postgres testcontainer: ", err)
	}

	if err := redisTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop redis testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestCheckDuplicates_DuplicatesNotExists(t *testing.T) {
	username := gofakeit.Username()
	result, err := credentialsRepository.CheckDuplicatesExists(context.Background(), username)

	require.NoError(t, err)
	assert.Equal(t, false, result)
}

func TestCheckDuplicates_DuplicatesExists(t *testing.T) {
	credentialsToSave := util.RandomDomainCredentials()
	err := credentialsRepository.Save(context.Background(), credentialsToSave)
	require.NoError(t, err)

	result, err := credentialsRepository.CheckDuplicatesExists(context.Background(), credentialsToSave.Username)

	require.NoError(t, err)
	assert.Equal(t, true, result)
}

func TestFindByUsername_UsernameNotFound(t *testing.T) {
	username := gofakeit.Username()
	result, err := credentialsRepository.FindByUsername(context.Background(), username)

	convertedError := status.Convert(err)
	assert.Equal(t, codes.NotFound, convertedError.Code())
	assert.Equal(t, credentials.ErrCredentialsNotFoundByUsername.Error(), convertedError.Message())

	assert.Nil(t, result)
}

func TestFindByUsername_UsernameFound(t *testing.T) {
	credentialsToSave := util.RandomDomainCredentials()
	err := credentialsRepository.Save(context.Background(), credentialsToSave)
	require.NoError(t, err)

	result, err := credentialsRepository.FindByUsername(context.Background(), credentialsToSave.Username)

	require.NoError(t, err)
	assert.Equal(t, *credentialsToSave, *result)
}

func TestSave_HappyPath(t *testing.T) {
	credentialsToSave := util.RandomDomainCredentials()
	foundCredentials, err := credentialsRepository.FindByUsername(context.Background(), credentialsToSave.Username)
	require.Error(t, err)
	assert.Nil(t, foundCredentials)

	err = credentialsRepository.Save(context.Background(), credentialsToSave)

	require.NoError(t, err)

	foundCredentials, err = credentialsRepository.FindByUsername(context.Background(), credentialsToSave.Username)
	require.NoError(t, err)
	assert.Equal(t, *credentialsToSave, *foundCredentials)
}
