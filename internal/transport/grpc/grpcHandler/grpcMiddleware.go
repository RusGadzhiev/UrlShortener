package grpcHandler

import (
	"context"

	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	"google.golang.org/grpc"
)

func LoggingUnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	logger.Info("New gRPC request: ", "method - ", info.FullMethod)

	m, err := handler(ctx, req)
	if err != nil {
		logger.Infof("RPC failed with error %v", err)
	}
	return m, err

}
