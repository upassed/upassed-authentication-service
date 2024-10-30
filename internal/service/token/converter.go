package token

import (
	"github.com/upassed/upassed-authentication-service/internal/jwt"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
)

func ConvertToBusinessTokenGenerateResponse(tokens *jwt.GeneratedTokens) *business.TokenGenerateResponse {
	return &business.TokenGenerateResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}
