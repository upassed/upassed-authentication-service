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
	require.NoError(t, err)

	assert.Equal(t, credentialsToConvert.ID, convertedCredentials.ID)
	assert.Equal(t, credentialsToConvert.Username, convertedCredentials.Username)
	assert.Equal(t, string(credentialsToConvert.AccountType), string(convertedCredentials.AccountType))

	err = bcrypt.CompareHashAndPassword(convertedCredentials.PasswordHash, []byte(credentialsToConvert.Password))
	require.NoError(t, err)
}
