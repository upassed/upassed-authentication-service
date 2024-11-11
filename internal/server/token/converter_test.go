package token_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/upassed/upassed-authentication-service/internal/server/token"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"testing"
)

func TestConvertToTokenGenerateRequest(t *testing.T) {
	requestToConvert := util.RandomClientTokenGenerateRequest()
	convertedRequest := token.ConvertToTokenGenerateRequest(requestToConvert)

	assert.Equal(t, requestToConvert.GetUsername(), convertedRequest.Username)
	assert.Equal(t, requestToConvert.GetPassword(), convertedRequest.Password)
}

func TestConvertToTokenGenerateResponse(t *testing.T) {
	responseToConvert := util.RandomBusinessTokenGenerateResponse()
	convertedResponse := token.ConvertToTokenGenerateResponse(responseToConvert)

	assert.Equal(t, responseToConvert.AccessToken, convertedResponse.GetAccessToken())
	assert.Equal(t, responseToConvert.RefreshToken, convertedResponse.GetRefreshToken())
}

func TestConvertToTokenRefreshRequest(t *testing.T) {
	requestToConvert := util.RandomClientTokenRefreshRequest()
	convertedRequest := token.ConvertToTokenRefreshRequest(requestToConvert)

	assert.Equal(t, requestToConvert.GetRefreshToken(), convertedRequest.RefreshToken)
}

func TestConvertToTokenRefreshResponse(t *testing.T) {
	responseToConvert := util.RandomBusinessTokenRefreshResponse()
	convertedResponse := token.ConvertToTokenRefreshResponse(responseToConvert)

	assert.Equal(t, responseToConvert.NewAccessToken, convertedResponse.GetNewAccessToken())
}

func TestConvertToTokenValidateRequest(t *testing.T) {
	requestToConvert := util.RandomClientTokenValidateRequest()
	convertedRequest := token.ConvertToTokenValidateRequest(requestToConvert)

	assert.Equal(t, requestToConvert.GetAccessToken(), convertedRequest.AccessToken)
}

func TestConvertToTokenAccessResponse(t *testing.T) {
	responseToConvert := util.RandomBusinessTokenValidateResponse()
	convertedResponse := token.ConvertToTokenValidateResponse(responseToConvert)

	assert.Equal(t, responseToConvert.Username, convertedResponse.GetUsername())
	assert.Equal(t, string(responseToConvert.AccountType), convertedResponse.GetAccountType())
	assert.Equal(t, responseToConvert.AccountID.String(), convertedResponse.GetAccountId())
}
