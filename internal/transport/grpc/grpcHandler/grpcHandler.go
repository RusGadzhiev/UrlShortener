package grpcHandler

import (
	"context"

	"github.com/RusGadzhiev/UrlShortener/internal/service"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
	"github.com/RusGadzhiev/UrlShortener/pkg/validator"
	proto "github.com/RusGadzhiev/UrlShortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service interface {
	// возвращает оригинальный url по короткому
	GetUrl(ctx context.Context, shortenUrl string) (string, error)
	// сокращает длинный урл, сохреняет его и возвращает укороченный урл
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
	shortUrl := r.GetShortUrl()

	if !validator.IsShortUrl(shortUrl) {
		logger.Debugf("ShortUrl: %s not valid", shortUrl)
		return nil, status.Error(codes.InvalidArgument, "")
	}

	longUrl, err := h.service.GetUrl(ctx, shortUrl)
	if err == service.ErrUrlNotFound {
		logger.Debugf("ShortUrl: %s not found", shortUrl)
		return nil, status.Error(codes.NotFound, "")
	} else if err != nil {
		logger.Errorf("ShortUrl: %s not found, err: %w", shortUrl, err)
		return nil, status.Error(codes.Internal, "")
	}

	return &proto.GetUrlResponse{LongUrl: longUrl}, nil

}

func (h *GRPCHandler) ShortenUrl(ctx context.Context, r *proto.ShortenUrlRequest) (*proto.ShortenUrlResponse, error) {
	longUrl := r.GetLongUrl()

	if !validator.IsUrl(longUrl) {
		logger.Debugf("LongUrl: %s not valid", longUrl)
		return nil, status.Error(codes.InvalidArgument, "")
	}

	shortUrl, err := h.service.ShortenUrl(ctx, longUrl)
	if err != nil {
		logger.Errorf("ShortenUrl err: %w", err)
		return nil, status.Error(codes.Internal, "")
	}

	return &proto.ShortenUrlResponse{ShortUrl: shortUrl}, nil

}


