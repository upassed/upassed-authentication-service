package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"reflect"
	"runtime"

	"github.com/upassed/upassed-authentication-service/internal/config"
	"github.com/upassed/upassed-authentication-service/internal/middleware"
	business "github.com/upassed/upassed-authentication-service/internal/service/model"
	"google.golang.org/grpc"
)

var (
	errStartingTcpConnection = errors.New("unable to start tcp connection")
	errStartingServer        = errors.New("unable to start gRPC server")
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

	return &AppServer{
		config: params.Config,
		log:    params.Log,
		server: server,
	}
}

func (server *AppServer) Run() error {
	op := runtime.FuncForPC(reflect.ValueOf(server.Run).Pointer()).Name()

	log := server.log.With(
		slog.String("op", op),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.GrpcServer.Port))
	if err != nil {
		return fmt.Errorf("%s -> %w; %w", op, errStartingTcpConnection, err)
	}

	log.Info("gRPC server is running", slog.String("address", listener.Addr().String()))
	if err := server.server.Serve(listener); err != nil {
		return fmt.Errorf("%s -> %w; %w", op, errStartingServer, err)
	}

	return nil
}

func (server *AppServer) GracefulStop() {
	op := runtime.FuncForPC(reflect.ValueOf(server.GracefulStop).Pointer()).Name()

	log := server.log.With(
		slog.String("op", op),
	)

	log.Info("gracefully stopping gRPC server...")
	server.server.GracefulStop()
}
