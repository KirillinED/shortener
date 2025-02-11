package storage

import (
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/storage/db"
	"github.com/KirillinED/shortener/internal/storage/interfaces"
)

func NewStorage(cfg *config.Config) (interfaces.Storage, error) {
	if cfg.DatabaseDSN != "" && cfg.DatabaseDriver != "" {
		return db.NewDbStorage(cfg)
	}

	memStore, err := NewMemoryStorage(cfg)
	if err != nil {
		return nil, err
	}

	return memStore, nil
}
