package token_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/upassed/upassed-authentication-service/internal/service/token"
	"github.com/upassed/upassed-authentication-service/internal/util"
	"testing"
)

func TestConvertToBusinessTokenGenerateResponse(t *testing.T) {
	tokens := util.RandomJwtGeneratedTokens()
	convertedResponse := token.ConvertToBusinessTokenGenerateResponse(tokens)

	assert.Equal(t, tokens.AccessToken, convertedResponse.AccessToken)
	assert.Equal(t, tokens.RefreshToken, convertedResponse.RefreshToken)
}
