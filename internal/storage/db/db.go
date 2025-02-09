package db

import (
	"errors"
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/storage/interfaces"
)

const PostgresDriver = "postgres"

func NewDbStorage(cfg *config.Config) (interfaces.Storage, error) {
	switch cfg.DatabaseDriver {
	case PostgresDriver:
		return NewPostgresStorage(cfg)
	}

	return nil, errors.New("db driver is not supported")
}
