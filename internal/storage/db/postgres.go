package db

import (
	"context"
	"errors"
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/entities"
	storageErrors "github.com/KirillinED/shortener/internal/storage/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresStorage struct {
	Pool *pgxpool.Pool
}

const (
	UniqueViolation = "23505"
)

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

func (ps *PostgresStorage) shortExists(url string) (bool, error) {
	ctx := context.Background()

	r := ps.Pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM urls WHERE long = $1)`, url)

	var res bool

	err := r.Scan(&res)
	if err != nil {
		return false, err
	}

	return res, nil
}

func (ps *PostgresStorage) longExists(url string) (bool, error) {
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

func (ps *PostgresStorage) StoreLink(link entities.Link) error {
	ctx := context.Background()

	_, err := ps.Pool.Exec(ctx, `INSERT INTO urls (short, long) VALUES ($1, $2)`, link.Short, link.Long)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && isDuplicateViolationError(pgErr) {
			return &storageErrors.DuplicateError{}
		}

		return err
	}

	return nil
}

func (ps *PostgresStorage) StoreLinks(links []entities.Link) error {
	batch := pgx.Batch{QueuedQueries: make([]*pgx.QueuedQuery, 0)}

	for _, link := range links {
		batch.Queue(`INSERT INTO urls (short, long) VALUES ($1, $2)`, link.Short, link.Long)
	}

	res := ps.Pool.SendBatch(context.Background(), &batch)
	err := res.Close()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && isDuplicateViolationError(pgErr) {
			return &storageErrors.DuplicateError{}
		}

		return err
	}

	return nil
}

func (ps *PostgresStorage) getUserLinks(id string) ([]entities.Link, error) {
	rows, err := ps.Pool.Query(
		context.Background(),
		`SELECT * FROM urls WHERE user_id = $1`,
		id,
	)

	if err != nil {
		return nil, err
	}

	var links []entities.Link
	err = rows.Scan(&links)
	if err != nil {
		return nil, err
	}

	return links, err
}

func (ps *PostgresStorage) Close() error {
	ps.Pool.Close()
	return nil
}

func isDuplicateViolationError(err *pgconn.PgError) bool {
	return err.Code == UniqueViolation
}
