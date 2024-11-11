package token

import (
	"github.com/upassed/upassed-authentication-service/internal/jwt"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
)

func ConvertToBusinessTokenGenerateResponse(tokens *jwt.GeneratedTokens) *business.TokenGenerateResponse {
	return &business.TokenGenerateResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}
}

func ConvertToBusinessTokenValidateResponse(credentials *domain.Credentials) *business.TokenValidateResponse {
	return &business.TokenValidateResponse{
		CredentialsID: credentials.ID,
		Username:      credentials.Username,
		AccountType:   business.AccountType(credentials.AccountType),
	}
}
