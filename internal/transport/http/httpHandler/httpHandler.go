package httpHandler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/RusGadzhiev/UrlShortener/internal/service"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	"github.com/RusGadzhiev/UrlShortener/pkg/validator"
	"github.com/gorilla/mux"
)

type Service interface {
	// возвращает оригинальный url по короткому
	GetUrl(ctx context.Context, shortenUrl string) (string, error)
	// сокращает длинный урл, сохреняет его и возвращает укороченный урл
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

func (h *HttpHandler) Router() *mux.Router {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/api/get-url", h.GetUrl).Methods("GET")
	r.HandleFunc("/api/shorten-url", h.ShortenUrl).Methods("POST")

	r.Use(func(hdl http.Handler) http.Handler {
		return h.PanicRecoverMiddleware(hdl)
	})
	r.Use(func(hdl http.Handler) http.Handler {
		return h.LoggingMiddleware(hdl)
	})
	
	return r
}

func (h *HttpHandler) GetUrl(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	shortUrl := r.URL.Query().Get("shortUrl")
	if !validator.IsShortUrl(shortUrl) {
		logger.Debugf("ShortUrl: %s not valid", shortUrl)
		h.clientError(w)
		return
	}

	longUrl, err := h.service.GetUrl(ctx, shortUrl)
	if err == service.ErrUrlNotFound {
		logger.Debugf("ShortUrl: %s not found", shortUrl)
		h.clientError(w)
		return
	} else if err != nil {
		logger.Errorf("ShortUrl: %s not found, err: %w", shortUrl, err)
		h.serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = renderJSON(w, longUrl)
	if err != nil {
		logger.Errorf("RenderJson err: %w", err)
	}
}

func (h *HttpHandler) ShortenUrl(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	longUrl := r.URL.Query().Get("longUrl")
	if !validator.IsUrl(longUrl) {
		logger.Debugf("LongUrl: %s not valid", longUrl)
		h.clientError(w)
		return
	}

	shortUrl, err := h.service.ShortenUrl(ctx, longUrl)
	if err != nil {
		logger.Errorf("ShortenUrl err: %w", err)
		h.serverError(w)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = renderJSON(w, shortUrl)
	if err != nil {
		logger.Errorf("RenderJson err: %w", err)
	}
}

// renderJSON преобразует 'v' в формат JSON и записывает результат, в виде ответа, в w.
func renderJSON(w http.ResponseWriter, v interface{}) error {
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
