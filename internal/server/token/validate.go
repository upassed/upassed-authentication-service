package token

import (
	"context"
	"github.com/upassed/upassed-authentication-service/internal/handling"
	requestid "github.com/upassed/upassed-authentication-service/internal/middleware/common/request_id"
	"github.com/upassed/upassed-authentication-service/internal/tracing"
	"github.com/upassed/upassed-authentication-service/pkg/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

func (server *tokenServerAPI) Validate(ctx context.Context, request *client.TokenValidateRequest) (*client.TokenValidateResponse, error) {
	spanContext, span := otel.Tracer(server.cfg.Tracing.TokenTracerName).Start(ctx, "token#Validate")
	span.SetAttributes(
		attribute.String(string(requestid.ContextKey), requestid.GetRequestIDFromContext(ctx)),
	)
	defer span.End()

	if err := request.Validate(); err != nil {
		tracing.SetSpanError(span, err)
		return nil, handling.Wrap(err, handling.WithCode(codes.InvalidArgument))
	}

	response, err := server.service.Validate(spanContext, ConvertToTokenValidateRequest(request))
	if err != nil {
		tracing.SetSpanError(span, err)
		return nil, err
	}

	return ConvertToTokenValidateResponse(response), nil
}
