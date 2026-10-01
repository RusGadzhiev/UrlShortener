package httpServer

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/RusGadzhiev/UrlShortener/internal/config"
	"github.com/RusGadzhiev/UrlShortener/internal/transport/http/httpHandler"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
)

type HttpServer struct {
	server http.Server
}

func NewHttpServer(h *httpHandler.HttpHandler, cfg config.Server) *HttpServer {
	return &HttpServer{
		server: http.Server{
			Addr:         ":" + cfg.Port,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
			Handler:      h.Router(),
		},
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	logger.Info("start listen http server", "addr", s.server.Addr)

	select {
	case <-ctx.Done():
		logger.Info("gracefully stopping", "cause", context.Cause(ctx))

		shtCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := s.server.Shutdown(shtCtx); err != nil {
			return err
		}
		return <-errCh
	case err := <-errCh:
		return err
	}
}
