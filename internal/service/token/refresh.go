package token

import (
	"context"
	"errors"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
)

func (service *tokenServiceImpl) Refresh(context.Context, *business.TokenRefreshRequest) (*business.TokenRefreshResponse, error) {
	return nil, errors.New("not implemented")
}
