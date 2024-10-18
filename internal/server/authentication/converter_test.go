package authentication_test

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/upassed/upassed-authentication-service/internal/server/authentication"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/pkg/client"
)

func TestConvertToCredentials(t *testing.T) {
	request := client.CredentialsCreateRequest{
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 40),
	}

	credentials := authentication.ConvertToCredentials(&request)

	assert.Equal(t, request.GetUsername(), credentials.Username)
	assert.Equal(t, request.GetPassword(), credentials.Password)
}

func TestConvertToCreateCredentialsResponse(t *testing.T) {
	credentialsResponseToConvert := business.CreateCredentialsResponse{
		CredentialsID: uuid.New(),
	}

	clientResponse := authentication.ConvertToCreateCredentialsResponse(credentialsResponseToConvert)

	assert.Equal(t, credentialsResponseToConvert.CredentialsID.String(), clientResponse.GetCredentialsId())
}
