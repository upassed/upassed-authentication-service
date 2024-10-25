package credentials

import (
	"context"
	"errors"
)

func (repository *credentialsRepositoryImpl) CheckDuplicatesExists(ctx context.Context, username string) (bool, error) {
	return false, errors.New("not implemented")
}
