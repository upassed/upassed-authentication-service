package token

import (
	"context"
	"github.com/upassed/upassed-authentication-service/internal/config"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"google.golang.org/grpc"
)

type tokenServerAPI struct {
	client.UnimplementedTokenServer
	cfg     *config.Config
	service tokenService
}

type tokenService interface {
	Generate(context.Context, *business.TokenGenerateRequest) (*business.TokenGenerateResponse, error)
	Refresh(context.Context, *business.TokenRefreshRequest) (*business.TokenRefreshResponse, error)
	Validate(ctx context.Context, request *business.TokenValidateRequest) (*business.TokenValidateResponse, error)
}

func Register(gRPC *grpc.Server, cfg *config.Config, service tokenService) {
	client.RegisterTokenServer(gRPC, &tokenServerAPI{
		cfg:     cfg,
		service: service,
	})
}
