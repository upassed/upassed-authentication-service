package credentials_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/service/credentials"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestConvertToDomainCredentials_HappyPath(t *testing.T) {
	credentialsToConvert := util.RandomBusinessCredentials()
	convertedCredentials, err := credentials.ConvertToDomainCredentials(credentialsToConvert)
	require.Nil(t, err)

	assert.Equal(t, credentialsToConvert.ID, convertedCredentials.ID)
	assert.Equal(t, credentialsToConvert.Username, convertedCredentials.Username)

	err = bcrypt.CompareHashAndPassword(convertedCredentials.PasswordHash, []byte(credentialsToConvert.Password))
	require.Nil(t, err)
}
