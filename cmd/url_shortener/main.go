package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/RusGadzhiev/UrlShortener/internal/config"
	"github.com/RusGadzhiev/UrlShortener/internal/service"
	"github.com/RusGadzhiev/UrlShortener/internal/storage/postgres"
	"github.com/RusGadzhiev/UrlShortener/internal/transport/http/httpHandler"
	"github.com/RusGadzhiev/UrlShortener/internal/transport/http/httpServer"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	"github.com/RusGadzhiev/UrlShortener/pkg/validator"
	"go.uber.org/zap"
)

type Server interface {
	Run(ctx context.Context, logger *zap.SugaredLogger) error
}

// нормально обработай ошибки
func main() {
	cfg := config.MustLoad()
	logger := logger.NewZapLogger()
	validator.ValidatorInit(cfg.Pattern)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	storage, err := postgres.NewPostgresStorage(ctx, cfg.PgDb)
	if err != nil {
		logger.Fatal(err)
	}
	logger.Info("Storage postgres started successfully")

	service := service.NewService(storage)

	httpHandler := httpHandler.NewHttpHandler(service, logger)

	server := httpServer.NewHttpServer(ctx, httpHandler, cfg.HTTPServer)

	if err := server.Run(ctx, logger); err != nil {
		logger.Fatal(err)
	}
}
