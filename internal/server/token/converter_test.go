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
