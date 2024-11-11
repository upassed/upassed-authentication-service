package token

import (
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/pkg/client"
)

func ConvertToTokenGenerateRequest(request *client.TokenGenerateRequest) *business.TokenGenerateRequest {
	return &business.TokenGenerateRequest{
		Username: request.GetUsername(),
		Password: request.GetPassword(),
	}
}

func ConvertToTokenGenerateResponse(response *business.TokenGenerateResponse) *client.TokenGenerateResponse {
	return &client.TokenGenerateResponse{
		AccessToken:  response.AccessToken,
		RefreshToken: response.RefreshToken,
	}
}

func ConvertToTokenRefreshRequest(request *client.TokenRefreshRequest) *business.TokenRefreshRequest {
	return &business.TokenRefreshRequest{
		RefreshToken: request.RefreshToken,
	}
}

func ConvertToTokenRefreshResponse(request *business.TokenRefreshResponse) *client.TokenRefreshResponse {
	return &client.TokenRefreshResponse{
		NewAccessToken: request.NewAccessToken,
	}
}

func ConvertToTokenValidateRequest(request *client.TokenValidateRequest) *business.TokenValidateRequest {
	return &business.TokenValidateRequest{
		AccessToken: request.AccessToken,
	}
}

func ConvertToTokenValidateResponse(request *business.TokenValidateResponse) *client.TokenValidateResponse {
	return &client.TokenValidateResponse{
		AccountId:   request.AccountID.String(),
		Username:    request.Username,
		AccountType: string(request.AccountType),
	}
}
