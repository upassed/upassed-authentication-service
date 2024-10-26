package credentials

import (
	"context"
	"errors"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
)

func (repository *credentialsRepositoryImpl) Save(ctx context.Context, credentials *domain.Credentials) error {
	return errors.New("not implemented")
}
