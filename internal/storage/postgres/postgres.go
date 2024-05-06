package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/RusGadzhiev/UrlShortener/internal/config"
	"github.com/RusGadzhiev/UrlShortener/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

var createQuery = `
	CREATE TABLE IF NOT EXISTS LINKS (
		short_url	VARCHAR(10) PRIMARY KEY,
		long_url		    VARCHAR(1024)
	);

	CREATE INDEX IF NOT EXISTS idx ON links USING hash(
		long_url
	);
`
var (
	ErrPingPostgres    = errors.New("error of ping postgres")
	ErrNewPoolPostgres = errors.New("error new pool postgres")
	ErrInitPostgres    = errors.New("error init postgres")
)

const (
	Short = "short"
	Long = "long"
)

type postgresStorage struct {
	pool *pgxpool.Pool
}

func NewPostgresStorage(ctx context.Context, cfg config.PgDb) (*postgresStorage, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, ErrNewPoolPostgres
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, ErrPingPostgres
	}

	_, err = pool.Exec(ctx, createQuery)
	if err != nil {
		return nil, ErrInitPostgres
	}

	return &postgresStorage{
		pool: pool,
	}, nil
}

func (s *postgresStorage) GetShortURL(ctx context.Context, longUrl string) (string, error) {
	return s.getURL(ctx, longUrl, Short)
}

func (s *postgresStorage) GetLongURL(ctx context.Context, shortUrl string) (string, error) {
	return s.getURL(ctx, shortUrl, Long)
}

func (s *postgresStorage) Add(ctx context.Context, longUrl string, shortUrl string) error {
	q := `INSERT INTO links(short_url, long_url) VALUES($1, $2)`

	_, err := s.pool.Exec(ctx, q, shortUrl, longUrl)
	return err
}

func (s *postgresStorage) getURL(ctx context.Context, url string, column string) (string, error) {
	var q string
	switch column {
		case Short:
			q = `SELECT short_url FROM links WHERE long_url = $1`
		case Long:
			q = `SELECT long_url FROM links WHERE short_url = $1`
	}

	var returnUrl string
	s.pool.QueryRow(ctx, q, url).Scan(&returnUrl)
	if returnUrl == "" {
		return "", service.ErrUrlNotFound
	}
	return returnUrl, nil
}