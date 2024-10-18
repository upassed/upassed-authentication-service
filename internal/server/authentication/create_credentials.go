package authentication

import (
	"context"

	"github.com/upassed/upassed-authentication-service/internal/handling"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"google.golang.org/grpc/codes"
)

func (server *authenticationServerAPI) CreateCredentials(ctx context.Context, request *client.CredentialsCreateRequest) (*client.CredentialsCreateResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	response, err := server.service.CreateCredentials(ctx, ConvertToCredentials(request))
	if err != nil {
		return nil, err
	}

	return ConvertToCreateCredentialsResponse(response), nil
}
