package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/RusGadzhiev/UrlShortener/internal/config"
	"github.com/RusGadzhiev/UrlShortener/internal/service"
	"github.com/RusGadzhiev/UrlShortener/internal/storage/postgres"
	"github.com/RusGadzhiev/UrlShortener/internal/transport/grpc/grpcHandler"
	"github.com/RusGadzhiev/UrlShortener/internal/transport/grpc/grpcServer"
	"github.com/RusGadzhiev/UrlShortener/internal/transport/http/httpHandler"
	"github.com/RusGadzhiev/UrlShortener/internal/transport/http/httpServer"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	"github.com/RusGadzhiev/UrlShortener/pkg/validator"
)

const (
	httpTransportMode = "http"
	grpcTransportMode = "grpc"
)

type Server interface {
	Run(ctx context.Context) error
}

func main() {
	cfg := config.MustLoad()
	validator.ValidatorInit(cfg.Pattern)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	storage, err := postgres.NewPostgresStorage(ctx, cfg.PgDb)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Storage postgres started successfully")

	service := service.NewService(storage)

	var server Server
	if cfg.TransportMode == httpTransportMode {
		httpHandler := httpHandler.NewHttpHandler(service)
		server = httpServer.NewHttpServer(ctx, httpHandler, cfg.Server)
	} else {
		grpcHandler := grpcHandler.NewGRPCHandler(service)
		server = grpcServer.NewGRPCServer(ctx, grpcHandler, cfg.Server.Port)
	}

	if err := server.Run(ctx); err != nil {
		logger.Fatal(err)
	}
}
