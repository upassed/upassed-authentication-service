package jwt_test

import (
	"github.com/upassed/upassed-authentication-service/internal/config"
	libjwt "github.com/upassed/upassed-authentication-service/internal/jwt"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

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

func TestGenerateTokens_HappyPath(t *testing.T) {
	username := gofakeit.Username()
	tokens, err := libjwt.New(cfg, logging.New(cfg.Env)).GenerateFor(username)
	assert.NoError(t, err)

	verifyToken(t, username, tokens.AccessToken, cfg.GetJwtAccessTokenTTL())
	verifyToken(t, username, tokens.RefreshToken, cfg.GetJwtRefreshTokenTTL())
}

func TestGenerateTokens_TokensExpires(t *testing.T) {
	oldAccessTokenTTL := cfg.Jwt.AccessTokenTTL
	cfg.Jwt.AccessTokenTTL = "-10m"

	username := gofakeit.Username()
	tokens, err := libjwt.New(cfg, logging.New(cfg.Env)).GenerateFor(username)
	assert.NoError(t, err)

	parsedToken, err := parseToken(tokens.AccessToken)
	assert.Error(t, err)
	assert.NotNil(t, parsedToken)
	assert.False(t, parsedToken.Valid)

	cfg.Jwt.AccessTokenTTL = oldAccessTokenTTL
}

func verifyToken(t *testing.T, username string, token string, tokenTTL time.Duration) {
	parsedToken, err := parseToken(token)

	assert.NoError(t, err)
	assert.NotNil(t, parsedToken)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, username, claims["username"])

	expirationDate, ok := claims["exp"].(float64)
	assert.True(t, ok)
	assert.WithinDuration(t, time.Unix(int64(expirationDate), 0), time.Now().Add(tokenTTL), 10*time.Second)
}

func parseToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKeyType
		}

		return []byte(cfg.Jwt.Secret), nil
	})
}
