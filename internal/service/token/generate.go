package token

import (
	"context"
	"errors"
	"github.com/upassed/upassed-authentication-service/internal/async"
	"github.com/upassed/upassed-authentication-service/internal/handling"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	domain "github.com/upassed/upassed-authentication-service/internal/repository/model"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
)

var (
	errFindingCredentialsByUsernameDeadlineExceeded = errors.New("finding credentials by username deadline exceeded")
	ErrPasswordHashNotMatch                         = errors.New("password hash does not match")
)

func (service *tokenServiceImpl) Generate(ctx context.Context, request *business.TokenGenerateRequest) (*business.TokenGenerateResponse, error) {
	spanContext, span := otel.Tracer(service.cfg.Tracing.TokenTracerName).Start(ctx, "tokenService#Generate")
	span.SetAttributes(attribute.String("username", request.Username))
	defer span.End()

	log := logging.Wrap(service.log,
		logging.WithOp(service.Generate),
		logging.WithCtx(ctx),
		logging.WithAny("username", request.Username),
	)

	log.Info("started generating tokens")
	timeout := service.cfg.GetEndpointExecutionTimeout()
	foundCredentials, err := async.ExecuteWithTimeout(spanContext, timeout, func(ctx context.Context) (*domain.Credentials, error) {
		log.Info("finding credentials by username")
		credentialsFromDatabase, err := service.credentialsRepository.FindByUsername(ctx, request.Username)
		if err != nil {
			log.Error("error while finding credentials by username")
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, err
		}

		log.Info("credentials found in database by username")
		return credentialsFromDatabase, nil
	})

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Error("credentials finding deadline exceeded")
			span.SetAttributes(attribute.String("err", err.Error()))
			return nil, handling.Wrap(errFindingCredentialsByUsernameDeadlineExceeded, handling.WithCode(codes.DeadlineExceeded))
		}

		log.Error("error while finding credentials", logging.Error(err))
		span.SetAttributes(attribute.String("err", err.Error()))
		return nil, handling.Process(err)
	}

	if err := bcrypt.CompareHashAndPassword(foundCredentials.PasswordHash, []byte(request.Password)); err != nil {
		log.Error("password does not match with hash")
		span.SetAttributes(attribute.String("err", ErrPasswordHashNotMatch.Error()))
		return nil, handling.Wrap(ErrPasswordHashNotMatch, handling.WithCode(codes.Internal))
	}

	tokens, err := service.tokenGenerator.GenerateFor(request.Username)
	if err != nil {
		log.Error("error while generating tokens", logging.Error(err))
		span.SetAttributes(attribute.String("err", ErrGeneratingTokens.Error()))
		return nil, handling.Wrap(ErrGeneratingTokens, handling.WithCode(codes.Internal))
	}

	log.Info("access and refresh tokens successfully generated")
	return ConvertToBusinessTokenGenerateResponse(tokens), nil
}
