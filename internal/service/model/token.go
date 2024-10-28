package business

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

type AccountType string

type TokenValidateResponse struct {
	Username    string
	AccountType AccountType
}
