package credentials_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/upassed/upassed-authentication-service/internal/messanging/credentials"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"testing"
)

func TestConvertToStudentCreateRequest_InvalidBytes(t *testing.T) {
	invalidBytes := make([]byte, 10)
	_, err := credentials.ConvertToCredentialsCreateRequest(invalidBytes)
	require.NotNil(t, err)
}

func TestConvertToStudentCreateRequest_ValidBytes(t *testing.T) {
	initialRequest := util.RandomEvenCredentialsCreateRequest()
	initialRequestBytes, err := json.Marshal(initialRequest)
	require.Nil(t, err)

	convertedRequest, err := credentials.ConvertToCredentialsCreateRequest(initialRequestBytes)
	require.Nil(t, err)

	assert.Equal(t, initialRequest, convertedRequest)
}

func TestConvertToStudent(t *testing.T) {
	request := util.RandomEvenCredentialsCreateRequest()
	convertedCredentials := credentials.ConvertToCredentials(request)
	require.NotNil(t, convertedCredentials)

	assert.NotNil(t, convertedCredentials.ID)
	assert.Equal(t, request.Username, convertedCredentials.Username)
	assert.Equal(t, request.Password, convertedCredentials.Password)
}
