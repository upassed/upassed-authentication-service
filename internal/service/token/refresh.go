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
	"log/slog"
)

var (
	errExtractingUsernameDeadlineExceeded = errors.New("extracting username from refresh token deadline exceeded")
)

func (service *tokenServiceImpl) Refresh(ctx context.Context, request *business.TokenRefreshRequest) (*business.TokenRefreshResponse, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.TokenTracerName).Start(ctx, "tokenService#Refresh")
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.Refresh),
		logging.WithCtx(ctx),
	)

	log.Info("started refreshing token")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	refreshedAccessToken, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (string, error) {
		parsedToken, err := service.parseToken(request.RefreshToken)
		if err != nil {
			log.Error("unable to parse refresh token", logging.Error(err))
			span.SetAttributes(attribute.String("err", ErrParsingToken.Error()))
			return "", ErrParsingToken
		}

		if !parsedToken.Valid {
			log.Error("refresh token is invalid")
			span.SetAttributes(attribute.String("err", ErrTokenInvalid.Error()))
			return "", ErrTokenInvalid
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			log.Error("unable to extract map claims from refresh token")
			span.SetAttributes(attribute.String("err", errExtractingTokenClaims.Error()))
			return "", errExtractingTokenClaims
		}

		username, ok := claims[libjwt.UsernameKey].(string)
		if !ok {
			log.Error("username key is not present in refresh token claims")
			span.SetAttributes(attribute.String("err", errUsernameClaimNotPresent.Error()))
			return "", errUsernameClaimNotPresent
		}

		log.Info("username successfully extracted", slog.String("username", username))
		tokens, err := service.tokenGenerator.GenerateFor(username)
		if err != nil {
			log.Error("error generating new tokens", logging.Error(err))
			span.SetAttributes(attribute.String("err", ErrGeneratingTokens.Error()))
			return "", ErrGeneratingTokens
		}

		log.Info("new access token generated")
		return tokens.AccessToken, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("extracting username deadline exceeded")
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, handling.Wrap(errExtractingUsernameDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while extracting username from refresh token", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, handling.Process(err)
	}

	log.Info("access token successfully refreshed")
	return &business.TokenRefreshResponse{
		NewAccessToken: refreshedAccessToken,
	}, nil
}
