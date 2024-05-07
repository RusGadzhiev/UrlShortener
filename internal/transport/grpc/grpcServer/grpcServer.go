package grpcServer

import (
	"context"
	"net"

	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	proto "github.com/RusGadzhiev/UrlShortener/proto"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	port   string
}

func NewGRPCServer(ctx context.Context, grpcHandlers proto.GRPCHandlerServer, port string) *GRPCServer {
	s := grpc.NewServer()
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
		if err := s.server.Serve(listener); err != nil && err != grpc.ErrServerStopped { // или http
			logger.Fatalf("Listen grpc server error: ", err)
		}
	}()
	logger.Info("Start listen grpc server at " + ":" + s.port)

	<-ctx.Done()
	logger.Info("Gracefully stopping...")

	s.server.GracefulStop()
	return err
}
