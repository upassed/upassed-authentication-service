package credentials

import (
	"encoding/json"
	"github.com/google/uuid"
	event "github.com/upassed/upassed-authentication-service/internal/messanging/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
)

func ConvertToCredentialsCreateRequest(messageBody []byte) (*event.CredentialsCreateRequest, error) {
	var request event.CredentialsCreateRequest
	if err := json.Unmarshal(messageBody, &request); err != nil {
		return nil, err
	}

	return &request, nil
}

func ConvertToCredentials(request *event.CredentialsCreateRequest) *business.Credentials {
	return &business.Credentials{
		ID:       uuid.New(),
		Username: request.Username,
		Password: request.Password,
	}
}
