package httpServer

import (
	"context"
	"net/http"
	"time"

	"github.com/RusGadzhiev/UrlShortener/internal/config"
	"github.com/RusGadzhiev/UrlShortener/internal/transport/http/httpHandler"
	"go.uber.org/zap"
)

type HttpServer struct {
	server http.Server
}

func NewHttpServer(ctx context.Context, h *httpHandler.HttpHandler, cfg config.HTTPServer) *HttpServer {
	return &HttpServer{
		server: http.Server{
			Addr: ":" + cfg.Port,
			ReadTimeout: cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout: cfg.IdleTimeout,
			Handler: h.Router(),
		},
	}
}

func (s *HttpServer) Run(ctx context.Context, logger *zap.SugaredLogger) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Listen error: ", err)
		}
	}()
	logger.Info("Start listen at " + s.server.Addr)

	<-ctx.Done()
	logger.Info("\n Gracefully stopping...")

	shtCtx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	err := s.server.Shutdown(shtCtx)
	return err
	
}