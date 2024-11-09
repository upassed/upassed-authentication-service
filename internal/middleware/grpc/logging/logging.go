package logging

import (
	"context"
	"fmt"
	logging "github.com/upassed/upassed-authentication-service/internal/logger"
	requestid "github.com/upassed/upassed-authentication-service/internal/middleware/common/request_id"
	"log/slog"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func MiddlewareInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		log := logging.Wrap(log, logging.WithOp(MiddlewareInterceptor))

		startTime := time.Now()
		resp, err := handler(ctx, req)
		elapsedTime := time.Since(startTime)

		log.Info("handled gRPC request",
			slog.String("requestID", requestid.GetRequestIDFromContext(ctx)),
			slog.String("method", info.FullMethod),
			slog.String("duration", fmt.Sprintf("%.2f ms", elapsedTime.Seconds()*1000)),
			slog.String("status", status.Code(err).String()),
		)

		return resp, err
	}
}
