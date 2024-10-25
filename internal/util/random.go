package util

import (
	"github.com/brianvoe/gofakeit/v7"
	event "github.com/upassed/upassed-authentication-service/internal/messanging/model"
)

func RandomEvenCredentialsCreateRequest() event.CredentialsCreateRequest {
	return event.CredentialsCreateRequest{
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 24),
	}
}
