package storage

import (
	"errors"
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/dto"
	"io"
)

type MemoryStorage struct {
	FileStorage         *FileStorage
	ShortToLongLinksMap map[string]string
	LongToShortLinksMap map[string]string
}

func NewMemoryStorage(cfg *config.Config) (*MemoryStorage, error) {
	memStore := &MemoryStorage{
		FileStorage:         NewFileStorage(cfg.FileStoragePath),
		ShortToLongLinksMap: make(map[string]string),
		LongToShortLinksMap: make(map[string]string),
	}

	err := memStore.Recovering()
	if err != nil {
		return nil, err
	}

	return memStore, nil
}

func (ms *MemoryStorage) ShortExists(shortURL string) (bool, error) {
	_, ok := ms.ShortToLongLinksMap[shortURL]

	return ok, nil
}

func (ms *MemoryStorage) LongExists(longURL string) (bool, error) {
	_, ok := ms.LongToShortLinksMap[longURL]

	return ok, nil
}

func (ms *MemoryStorage) GetShortURL(longURL string) (string, error) {
	if ok, _ := ms.LongExists(longURL); ok {
		return ms.LongToShortLinksMap[longURL], nil
	}

	return "", nil
}

func (ms *MemoryStorage) GetLongURL(shortURL string) (string, error) {
	if ok, _ := ms.ShortExists(shortURL); ok {
		return ms.ShortToLongLinksMap[shortURL], nil
	}

	return "", nil
}

func (ms *MemoryStorage) StoreLink(link dto.Link) (bool, error) {
	err := ms.FileStorage.encoder.Encode(link)
	if err != nil {
		return false, err
	}

	ms.ShortToLongLinksMap[link.Short] = link.Long

	ms.LongToShortLinksMap[link.Long] = link.Short

	return true, nil
}

func (ms *MemoryStorage) Close() error {
	err := ms.FileStorage.Close()
	if err != nil {
		return err
	}

	return nil
}

func (ms *MemoryStorage) Recovering() error {
	for {
		link, err := ms.FileStorage.ReadLink()
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				break
			}
			return err
		}

		_, err = ms.StoreLink(*link)
		if err != nil {
			return err
		}
	}

	return nil
}
