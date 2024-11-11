package business

import "github.com/google/uuid"

type TokenGenerateRequest struct {
	Username string
	Password string
}

type TokenGenerateResponse struct {
	AccessToken  string
	RefreshToken string
}

type TokenRefreshRequest struct {
	RefreshToken string
}

type TokenRefreshResponse struct {
	NewAccessToken string
}

type TokenValidateRequest struct {
	AccessToken string
}

type TokenValidateResponse struct {
	CredentialsID uuid.UUID
	Username      string
	AccountType   AccountType
}
