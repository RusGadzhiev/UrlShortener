package grpcServer

import (
	"context"
	"errors"
	"net"

	"github.com/RusGadzhiev/UrlShortener/internal/transport/grpc/grpcHandler"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	proto "github.com/RusGadzhiev/UrlShortener/proto"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	port   string
}

func NewGRPCServer(grpcHandlers proto.GRPCHandlerServer, port string) *GRPCServer {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpcHandler.LoggingUnaryServerInterceptor),
	)
	proto.RegisterGRPCHandlerServer(s, grpcHandlers)

	return &GRPCServer{
		server: s,
		port:   port,
	}
}

func (s *GRPCServer) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}

	go func() {
		if err := s.server.Serve(listener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			logger.Fatal("listen grpc server", "err", err)
		}
	}()
	logger.Info("start listen grpc server", "addr", ":"+s.port)

	<-ctx.Done()
	logger.Info("gracefully stopping", "cause", context.Cause(ctx))

	s.server.GracefulStop()
	return err
}
