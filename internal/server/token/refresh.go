package token

import (
	"context"
	"github.com/upassed/upassed-authentication-service/internal/handling"
	"github.com/upassed/upassed-authentication-service/internal/middleware"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

func (server *tokenServerAPI) Refresh(ctx context.Context, request *client.TokenRefreshRequest) (*client.TokenRefreshResponse, error) {
	spanContext, span := otel.Tracer(server.cfg.Tracing.TokenTracerName).Start(ctx, "token#Refresh")
	span.SetAttributes(
		attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
	)
	defer span.End()

	if err := request.Validate(); err != nil {
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	response, err := server.service.Refresh(spanContext, ConvertToTokenRefreshRequest(request))
	if err != nil {
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, err
	}

	return ConvertToTokenRefreshResponse(response), nil
}
