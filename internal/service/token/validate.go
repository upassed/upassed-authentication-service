package token

import (
	"context"
	"errors"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
)

func (service *tokenServiceImpl) Validate(context.Context, *business.TokenValidateRequest) (*business.TokenValidateResponse, error) {
	return nil, errors.New("not implemented")
}
