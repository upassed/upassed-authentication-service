package credentials_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/caching"
	"github.com/upassed/upassed-authentication-service/internal/caching/credentials"
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/testcontainer"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

var (
	redisClient *credentials.RedisClient
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
	logger := logging.New(cfg.Env)

	redisTestcontainer, err := testcontainer.NewRedisTestcontainer(ctx, cfg)
	if err != nil {
		log.Fatal("unable to run redis testcontainer: ", err)
	}

	port, err := redisTestcontainer.Start(ctx)
	if err != nil {
		log.Fatal("unable to get a postgres testcontainer real port: ", err)
	}

	cfg.Redis.Port = strconv.Itoa(port)
	redis, err := caching.OpenRedisConnection(cfg, logger)
	if err != nil {
		log.Fatal("unable to open connections to redis: ", err)
	}

	redisClient = credentials.New(redis, cfg, logger)
	exitCode := m.Run()
	if err := redisTestcontainer.Stop(ctx); err != nil {
		log.Println("unable to stop redis testcontainer: ", err)
	}

	os.Exit(exitCode)
}

func TestSaveGroup_HappyPath(t *testing.T) {
	credentialsToSave := util.RandomDomainCredentials()
	ctx := context.Background()
	err := redisClient.Save(ctx, credentialsToSave)
	require.Nil(t, err)

	credentialsFromCache, err := redisClient.GetByID(ctx, credentialsToSave.ID)
	require.Nil(t, err)

	assert.Equal(t, *credentialsToSave, *credentialsFromCache)
}

func TestFindGroupByID_GroupNotFound(t *testing.T) {
	credentialsID := uuid.New()
	foundCredentials, err := redisClient.GetByID(context.Background(), credentialsID)
	require.NotNil(t, err)

	assert.ErrorIs(t, err, credentials.ErrCredentialsIsNotPresentInCache)
	assert.Nil(t, foundCredentials)
}
