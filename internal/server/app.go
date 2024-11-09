package server

import (
	"errors"
	"fmt"
	"github.com/upassed/upassed-authentication-service/internal/config"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	loggingMW "github.com/upassed/upassed-authentication-service/internal/middleware/grpc/logging"
	"github.com/upassed/upassed-authentication-service/internal/middleware/grpc/recovery"
	"github.com/upassed/upassed-authentication-service/internal/middleware/grpc/requestid"
	"github.com/upassed/upassed-authentication-service/internal/server/token"
	tokenSvc "github.com/upassed/upassed-authentication-service/internal/service/token"
	"google.golang.org/grpc"
	"log/slog"
	"net"
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
	Config       *config.Config
	Log          *slog.Logger
	TokenService tokenSvc.Service
}

func New(params AppServerCreateParams) *AppServer {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestid.MiddlewareInterceptor(),
			recovery.MiddlewareInterceptor(params.Log),
			loggingMW.MiddlewareInterceptor(params.Log),
		),
	)

	token.Register(server, params.Config, params.TokenService)
	return &AppServer{
		config: params.Config,
		log:    params.Log,
		server: server,
	}
}

func (server *AppServer) Run() error {
	log := logging.Wrap(server.log,
		logging.WithOp(server.Run),
		logging.WithAny("port", server.config.GrpcServer.Port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", server.config.GrpcServer.Port))
	if err != nil {
		log.Error("error while starting tcp connection", logging.Error(err))
		return errStartingTcpConnection
	}

	log.Info("gRPC server is running", slog.String("address", listener.Addr().String()))
	if err := server.server.Serve(listener); err != nil {
		log.Error("error while starting tcp server", logging.Error(err))
		return errStartingServer
	}

	log.Info("tcp server is now running")
	return nil
}

func (server *AppServer) GracefulStop() {
	log := logging.Wrap(server.log, logging.WithOp(server.GracefulStop))

	log.Info("gracefully stopping gRPC server...")
	server.server.GracefulStop()
}
