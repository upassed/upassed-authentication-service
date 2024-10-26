package util

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	event "github.com/upassed/upassed-authentication-service/internal/messanging/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
)

func RandomEvenCredentialsCreateRequest() *event.CredentialsCreateRequest {
	return &event.CredentialsCreateRequest{
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 24),
	}
}

func RandomBusinessCredentials() *business.Credentials {
	return &business.Credentials{
		ID:       uuid.New(),
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 24),
	}
}
