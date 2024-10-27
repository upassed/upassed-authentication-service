package util

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	event "github.com/upassed/upassed-authentication-service/internal/messanging/model"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"golang.org/x/crypto/bcrypt"
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

func RandomDomainCredentials() *domain.Credentials {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(gofakeit.Password(true, true, true, true, true, 24)), bcrypt.DefaultCost)

	return &domain.Credentials{
		ID:           uuid.New(),
		Username:     gofakeit.Username(),
		PasswordHash: passwordHash,
	}
}
