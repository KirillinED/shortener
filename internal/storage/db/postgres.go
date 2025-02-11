package db

import (
	"context"
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	Pool *pgxpool.Pool
}

func NewPostgresStorage(cfg *config.Config) (*PostgresStorage, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{Pool: pool}, nil
}

func (ps *PostgresStorage) ShortExists(url string) (bool, error) {
	ctx := context.Background()

	r := ps.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM urls WHERE long = $1)`, url)

	var res bool

	err := r.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (ps *PostgresStorage) LongExists(url string) (bool, error) {
	ctx := context.Background()

	r := ps.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM urls WHERE long = $1)`, url)

	var res bool

	err := r.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (ps *PostgresStorage) GetShortURL(url string) (string, error) {
	ctx := context.Background()

	r := ps.Pool.QueryRow(ctx, `SELECT short FROM urls WHERE long = $1`, url)

	var res string

	err := r.Scan(&res)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (ps *PostgresStorage) GetLongURL(url string) (string, error) {
	ctx := context.Background()

	r := ps.Pool.QueryRow(ctx, `SELECT long FROM urls WHERE short = $1`, url)

	var res string

	err := r.Scan(&res)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (ps *PostgresStorage) StoreLink(link dto.Link) (bool, error) {
	ctx := context.Background()

	c, err := ps.Pool.Exec(ctx, `INSERT INTO urls (short, long) VALUES ($1, $2)`, link.Short, link.Long)
	if err != nil {
		return false, err
	}

	return c.Insert(), nil
}

func (ps *PostgresStorage) Close() error {
	ps.Pool.Close()
	return nil
}
