package authentication

import (
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/pkg/client"
)

func ConvertToCredentials(request *client.CredentialsCreateRequest) business.Credentials {
	return business.Credentials{
		Username: request.GetUsername(),
		Password: request.GetPassword(),
	}
}

func ConvertToCreateCredentialsResponse(credentialsToConver business.CreateCredentialsResponse) *client.CredentialsCreateResponse {
	return &client.CredentialsCreateResponse{
		CredentialsId: credentialsToConver.CredentialsID.String(),
	}
}
