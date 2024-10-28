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

func (server *tokenServerAPI) Generate(ctx context.Context, request *client.TokenGenerateRequest) (*client.TokenGenerateResponse, error) {
	spanContext, span := otel.Tracer(server.cfg.Tracing.TokenTracerName).Start(ctx, "token#Generate")
	span.SetAttributes(
		attribute.String(string(middleware.RequestIDKey), middleware.GetRequestIDFromContext(ctx)),
		attribute.String("username", request.GetUsername()),
	)
	defer span.End()

	if err := request.Validate(); err != nil {
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	response, err := server.service.Generate(spanContext, ConvertToTokenGenerateRequest(request))
	if err != nil {
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, err
	}

	return ConvertToTokenGenerateResponse(response), nil
}
