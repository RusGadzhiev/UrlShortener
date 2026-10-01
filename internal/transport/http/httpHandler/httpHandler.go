package httpHandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/RusGadzhiev/UrlShortener/internal/service"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	"github.com/RusGadzhiev/UrlShortener/pkg/validator"
)

type Service interface {
	GetUrl(ctx context.Context, shortenURL string) (string, error)
	ShortenUrl(ctx context.Context, url string) (string, error)
}

type HttpHandler struct {
	service Service
}

func NewHttpHandler(service Service) *HttpHandler {
	return &HttpHandler{
		service: service,
	}
}

func (h *HttpHandler) Router() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/get-url", h.GetUrl)
	mux.HandleFunc("POST /api/shorten-url", h.ShortenUrl)

	return h.PanicRecoverMiddleware(h.LoggingMiddleware(mux))
}

func (h *HttpHandler) GetUrl(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	shortURL := r.URL.Query().Get("shortUrl")
	if !validator.IsShortUrl(shortURL) {
		logger.Debug("short url is not valid", "short_url", shortURL)
		h.clientError(w)
		return
	}

	longURL, err := h.service.GetUrl(ctx, shortURL)
	if errors.Is(err, service.ErrUrlNotFound) {
		logger.Debug("short url not found", "short_url", shortURL)
		h.clientError(w)
		return
	}
	if err != nil {
		logger.Error("get url", "short_url", shortURL, "err", err)
		h.serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = renderJSON(w, longURL)
	if err != nil {
		logger.Error("render json", "err", err)
	}
}

func (h *HttpHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	longURL := r.URL.Query().Get("longUrl")
	if !validator.IsUrl(longURL) {
		logger.Debug("long url is not valid", "long_url", longURL)
		h.clientError(w)
		return
	}

	shortURL, err := h.service.ShortenUrl(ctx, longURL)
	if err != nil {
		logger.Error("shorten url", "err", err)
		h.serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = renderJSON(w, shortURL)
	if err != nil {
		logger.Error("render json", "err", err)
	}
}

// renderJSON преобразует 'v' в формат JSON и записывает результат, в виде ответа, в w.
func renderJSON(w http.ResponseWriter, v any) error {
	json, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(json)
	return err
}

func (h *HttpHandler) serverError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (h *HttpHandler) clientError(w http.ResponseWriter) {
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
