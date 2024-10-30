package token

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt"
	"github.com/upassed/upassed-authentication-service/internal/async"
	"github.com/upassed/upassed-authentication-service/internal/handling"
	libjwt "github.com/upassed/upassed-authentication-service/internal/jwt"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc/codes"
)

var (
	errValidatingTokenDeadlineExceeded = errors.New("token validation deadline exceeded")
)

func (service *tokenServiceImpl) Validate(ctx context.Context, request *business.TokenValidateRequest) (*business.TokenValidateResponse, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.TokenTracerName).Start(ctx, "tokenService#Validate")
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.Validate),
		logging.WithCtx(ctx),
	)

	log.Info("started validating token")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	response, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*business.TokenValidateResponse, error) {
		parsedToken, err := service.parseToken(request.AccessToken)
		if err != nil {
			log.Error("unable to parse refresh token", logging.Error(err))
			span.SetAttributes(attribute.String("err", ErrParsingToken.Error()))
			return nil, ErrParsingToken
		}

		if !parsedToken.Valid {
			log.Error("refresh token is invalid")
			span.SetAttributes(attribute.String("err", ErrTokenInvalid.Error()))
			return nil, ErrTokenInvalid
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			log.Error("unable to extract map claims from refresh token")
			span.SetAttributes(attribute.String("err", errExtractingTokenClaims.Error()))
			return nil, errExtractingTokenClaims
		}

		username, ok := claims[libjwt.UsernameKey].(string)
		if !ok {
			log.Error("username key is not present in refresh token claims")
			span.SetAttributes(attribute.String("err", errUsernameClaimNotPresent.Error()))
			return nil, errUsernameClaimNotPresent
		}

		foundCredentials, err := service.credentialsRepository.FindByUsername(ctx, username)
		if err != nil {
			log.Error("error while finding credentials by username")
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, err
		}

		log.Info("access token validated successfully")
		return &business.TokenValidateResponse{
			Username:    username,
			AccountType: business.AccountType(foundCredentials.AccountType),
		}, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("validating token deadline exceeded")
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, handling.Wrap(errValidatingTokenDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while validating access token", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, handling.Process(err)
	}

	log.Info("token validation was successful")
	return response, nil
}
