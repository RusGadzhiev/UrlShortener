package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RusGadzhiev/UrlShortener/internal/service/encoder"
	"github.com/RusGadzhiev/UrlShortener/internal/storage/postgres"
	"github.com/RusGadzhiev/UrlShortener/pkg/logger"
)

var (
	ErrUrlNotFound = errors.New("no such url")
)

type Storage interface {
	GetShortURL(ctx context.Context, longURL string) (string, error)
	GetLongURL(ctx context.Context, shortURL string) (string, error)
	Add(ctx context.Context, longURL string, shortURL string) error
}

type service struct {
	repo Storage
}

func NewService(repo Storage) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) GetUrl(ctx context.Context, shortURL string) (string, error) {
	longURL, err := s.repo.GetLongURL(ctx, shortURL)
	if errors.Is(err, postgres.ErrNotFound) {
		logger.Debug("url not found", "short_url", shortURL, "err", err)
		return "", ErrUrlNotFound
	}
	if err != nil {
		return "", fmt.Errorf("GetUrl error: %w", err)
	}
	return longURL, nil
}

func (s *service) ShortenUrl(ctx context.Context, longURL string) (string, error) {
	shortURL, err := s.repo.GetShortURL(ctx, longURL)
	if err == nil {
		logger.Debug("short url already exists", "short_url", shortURL)
		return shortURL, nil
	}
	if !errors.Is(err, postgres.ErrNotFound) {
		return "", fmt.Errorf("ShortenUrl (GetShortURL) error: %w", err)
	}

	shortURL = encoder.Encode(int(time.Now().Unix()))

	err = s.repo.Add(ctx, longURL, shortURL)
	if err != nil {
		return "", fmt.Errorf("ShortenUrl (Add) error: %w", err)
	}

	return shortURL, nil
}
