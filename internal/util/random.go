package util

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/upassed/upassed-authentication-service/internal/jwt"
	event "github.com/upassed/upassed-authentication-service/internal/messanging/model"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"golang.org/x/crypto/bcrypt"
)

func RandomEvenCredentialsCreateRequest() *event.CredentialsCreateRequest {
	return &event.CredentialsCreateRequest{
		Username:    gofakeit.Username(),
		Password:    gofakeit.Password(true, true, true, true, true, 24),
		AccountType: event.TeacherAccountType,
	}
}

func RandomBusinessCredentials() *business.Credentials {
	return &business.Credentials{
		ID:          uuid.New(),
		Username:    gofakeit.Username(),
		Password:    gofakeit.Password(true, true, true, true, true, 24),
		AccountType: business.TeacherAccountType,
	}
}

func RandomDomainCredentials() *domain.Credentials {
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(gofakeit.Password(true, true, true, true, true, 24)), bcrypt.DefaultCost)

	return &domain.Credentials{
		ID:           uuid.New(),
		Username:     gofakeit.Username(),
		PasswordHash: passwordHash,
		AccountType:  domain.TeacherAccountType,
	}
}

func RandomClientTokenGenerateRequest() *client.TokenGenerateRequest {
	return &client.TokenGenerateRequest{
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 24),
	}
}

func RandomClientTokenRefreshRequest() *client.TokenRefreshRequest {
	return &client.TokenRefreshRequest{
		RefreshToken: gofakeit.Slogan(),
	}
}

func RandomClientTokenValidateRequest() *client.TokenValidateRequest {
	return &client.TokenValidateRequest{
		AccessToken: gofakeit.Slogan(),
	}
}

func RandomBusinessTokenGenerateRequest() *business.TokenGenerateRequest {
	return &business.TokenGenerateRequest{
		Username: gofakeit.Username(),
		Password: gofakeit.Password(true, true, true, true, true, 24),
	}
}

func RandomBusinessTokenGenerateResponse() *business.TokenGenerateResponse {
	return &business.TokenGenerateResponse{
		AccessToken:  gofakeit.Slogan(),
		RefreshToken: gofakeit.Slogan(),
	}
}

func RandomBusinessTokenRefreshResponse() *business.TokenRefreshResponse {
	return &business.TokenRefreshResponse{
		NewAccessToken: gofakeit.Slogan(),
	}
}

func RandomBusinessTokenValidateResponse() *business.TokenValidateResponse {
	return &business.TokenValidateResponse{
		CredentialsID: uuid.New(),
		Username:      gofakeit.Username(),
		AccountType:   business.StudentAccountType,
	}
}

func RandomBusinessTokenRefreshRequest() *business.TokenRefreshRequest {
	return &business.TokenRefreshRequest{
		RefreshToken: gofakeit.Slogan(),
	}
}

func RandomBusinessTokenValidateRequest() *business.TokenValidateRequest {
	return &business.TokenValidateRequest{
		AccessToken: gofakeit.Slogan(),
	}
}

func RandomJwtGeneratedTokens() *jwt.GeneratedTokens {
	return &jwt.GeneratedTokens{
		AccessToken:  gofakeit.Slogan(),
		RefreshToken: gofakeit.Slogan(),
	}
}
