package token

import (
	"context"
	"errors"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
)

func (service *tokenServiceImpl) Generate(context.Context, *business.TokenGenerateRequest) (*business.TokenGenerateResponse, error) {
	return nil, errors.New("not implemented")
}
