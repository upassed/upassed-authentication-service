package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/middleware"
	"github.com/upassed/upassed-authentication-service/internal/server/authentication"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"google.golang.org/grpc"
)

var (
	ErroStartingTcpConnection error = errors.New("unable to start tcp connection")
	ErrStartingServer         error = errors.New("unable to start gRPC server")
)

type AppServer struct {
	config *config.Config
	log    *slog.Logger
	server *grpc.Server
}

type AppServerCreateParams struct {
	Config                *config.Config
	Log                   *slog.Logger
	AuthenticationService authenticationService
}

type authenticationService interface {
	CreateCredentials(context.Context, business.Credentials) (business.CreateCredentialsResponse, error)
}

func New(params AppServerCreateParams) *AppServer {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.RequestIDMiddlewareInterceptor(),
			middleware.PanicRecoveryMiddlewareInterceptor(params.Log),
			middleware.LoggingMiddlewareInterceptor(params.Log),
		),
	)

	authentication.Register(server, params.AuthenticationService)
	return &AppServer{
		config: params.Config,
		log:    params.Log,
		server: server,
	}
}

func (server *AppServer) Run() error {
	const op = "server.Run()"

	log := server.log.With(
		slog.String("op", op),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.GrpcServer.Port))
	if err != nil {
		return fmt.Errorf("%s -> %w; %w", op, ErroStartingTcpConnection, err)
	}

	log.Info("gRPC server is running", slog.String("address", listener.Addr().String()))
	if err := server.server.Serve(listener); err != nil {
		return fmt.Errorf("%s -> %w; %w", op, ErrStartingServer, err)
	}

	return nil
}

func (server *AppServer) GracefulStop() {
	const op = "server.GracefulStop()"

	log := server.log.With(
		slog.String("op", op),
	)

	log.Info("gracefully stopping gRPC server...")
	server.server.GracefulStop()
}
