package grpcHandler

import (
	"context"
	"errors"

	"github.com/RusGadzhiev/UrlShortener/internal/service"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	"github.com/RusGadzhiev/UrlShortener/pkg/validator"
	proto "github.com/RusGadzhiev/UrlShortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	GetUrl(ctx context.Context, shortenURL string) (string, error)
	ShortenUrl(ctx context.Context, url string) (string, error)
}

type GRPCHandler struct {
	service Service
	proto.UnimplementedGRPCHandlerServer
}

func NewGRPCHandler(s Service) *GRPCHandler {
	return &GRPCHandler{
		service: s,
	}
}

func (h *GRPCHandler) GetUrl(ctx context.Context, r *proto.GetUrlRequest) (*proto.GetUrlResponse, error) {
	shortURL := r.GetShortUrl()

	if !validator.IsShortUrl(shortURL) {
		logger.Debug("short url is not valid", "short_url", shortURL)
		return nil, status.Error(codes.InvalidArgument, "")
	}

	longURL, err := h.service.GetUrl(ctx, shortURL)
	if errors.Is(err, service.ErrUrlNotFound) {
		logger.Debug("short url not found", "short_url", shortURL)
		return nil, status.Error(codes.NotFound, "")
	}
	if err != nil {
		logger.Error("get url", "short_url", shortURL, "err", err)
		return nil, status.Error(codes.Internal, "")
	}

	return &proto.GetUrlResponse{LongUrl: longURL}, nil
}

func (h *GRPCHandler) ShortenUrl(ctx context.Context, r *proto.ShortenUrlRequest) (*proto.ShortenUrlResponse, error) {
	longURL := r.GetLongUrl()

	if !validator.IsUrl(longURL) {
		logger.Debug("long url is not valid", "long_url", longURL)
		return nil, status.Error(codes.InvalidArgument, "")
	}

	shortURL, err := h.service.ShortenUrl(ctx, longURL)
	if err != nil {
		logger.Error("shorten url", "err", err)
		return nil, status.Error(codes.Internal, "")
	}

	return &proto.ShortenUrlResponse{ShortUrl: shortURL}, nil
}
