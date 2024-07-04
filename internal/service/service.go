package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RusGadzhiev/UrlShortener/internal/service/encoder"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
)

var (
	ErrUrlNotFound = errors.New("no such url")
)

type Storage interface {
	GetShortURL(ctx context.Context, longUrl string) (string, error)
	GetLongURL(ctx context.Context, shortUrl string) (string, error)
	Add(ctx context.Context, longUrl string, shortUrl string) error
}

type service struct {
	repo Storage
}

func NewService(repo Storage) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetUrl(ctx context.Context, shortUrl string) (string, error) {
	longUrl, err := s.repo.GetLongURL(ctx, shortUrl)
	if err == ErrUrlNotFound {
		logger.Debugf("LongUrl: %s  err: %w", longUrl, err)
		return "", ErrUrlNotFound
	} else if err != nil {
		return "", fmt.Errorf("GetUrl error: %w", err)
	}
	return longUrl, err
}

func (s *service) ShortenUrl(ctx context.Context, longUrl string) (string, error) {
	shortUrl, err := s.repo.GetShortURL(ctx, longUrl)
	if err == nil {
		logger.Debugf("ShortUrl: %s already exist", shortUrl)
		return shortUrl, nil
	} else if err != ErrUrlNotFound {
		return "", fmt.Errorf("ShortenUrl (GetShortURL) error: %w", err)
	}

	shortUrl = encoder.Encode(int(time.Now().Unix()))

	err = s.repo.Add(ctx, longUrl, shortUrl)
	if err != nil {
		return "", fmt.Errorf("ShortenUrl (Add) error: %w", err)
	}

	return shortUrl, nil
}
