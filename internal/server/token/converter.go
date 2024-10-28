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
