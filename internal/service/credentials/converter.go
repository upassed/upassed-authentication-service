package credentials

import (
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"golang.org/x/crypto/bcrypt"
)

func ConvertToDomainCredentials(credentials *business.Credentials) (*domain.Credentials, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &domain.Credentials{
		ID:           credentials.ID,
		Username:     credentials.Username,
		PasswordHash: passwordHash,
	}, nil
}
