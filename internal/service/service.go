package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/RusGadzhiev/UrlShortener/internal/service/encoder"
)

var (
	ErrUrlNotFound = errors.New("no such url")
	ErrUrlExist = errors.New("url exist")
)

type Storage interface {
	// возвращает ErrUrlNotFound если урла нет
	GetShortURL(ctx context.Context, longUrl string) (string, error)
	// возвращает ErrUrlNotFound если урла нет
	GetLongURL(ctx context.Context, shortUrl string) (string, error)
	// добавляет урл
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
		return "", fmt.Errorf("GetUrl error: %w", err)
	} else if err != nil {
		return "", fmt.Errorf("GetUrl error: %w", err)
	}
	return longUrl, err
}

// сделать с минимумом обращений в базу
func (s *service) ShortenUrl(ctx context.Context, longUrl string) (string, error) {
	// длинная ссылка уже сохранена
	if shortUrl, err := s.repo.GetShortURL(ctx, longUrl); err != ErrUrlNotFound {
		return shortUrl, nil
	}

	shortUrl := encoder.Encode(int(time.Now().Unix()))

	err := s.repo.Add(ctx, longUrl, shortUrl)
	if err != nil {
		return "", err
	}

	return shortUrl, nil
}
