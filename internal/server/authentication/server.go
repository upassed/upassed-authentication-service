package authentication

import (
	"context"

	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"google.golang.org/grpc"
)

type authenticationServerAPI struct {
	client.UnimplementedAuthenticationServer
	service authenticationService
}

type authenticationService interface {
	CreateCredentials(context.Context, business.Credentials) (business.CreateCredentialsResponse, error)
}

func Register(gRPC *grpc.Server, service authenticationService) {
	client.RegisterAuthenticationServer(gRPC, &authenticationServerAPI{
		service: service,
	})
}
