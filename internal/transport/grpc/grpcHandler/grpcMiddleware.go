package grpcHandler

import (
	"context"

	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	"google.golang.org/grpc"
)

func LoggingUnaryServerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	logger.Info("new grpc request", "method", info.FullMethod)

	m, err := handler(ctx, req)
	if err != nil {
		logger.Info("rpc failed", "err", err)
	}
	return m, err
}
