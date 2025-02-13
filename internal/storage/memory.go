package storage

import (
	"errors"
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/dto"
	storageErrors "github.com/KirillinED/shortener/internal/storage/errors"
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

func (ms *MemoryStorage) shortExists(shortURL string) (bool, error) {
	_, ok := ms.ShortToLongLinksMap[shortURL]

	return ok, nil
}

func (ms *MemoryStorage) longExists(longURL string) (bool, error) {
	_, ok := ms.LongToShortLinksMap[longURL]

	return ok, nil
}

func (ms *MemoryStorage) GetShortURL(longURL string) (string, error) {
	if ok, _ := ms.longExists(longURL); ok {
		return ms.LongToShortLinksMap[longURL], nil
	}

	return "", nil
}

func (ms *MemoryStorage) GetLongURL(shortURL string) (string, error) {
	if ok, _ := ms.shortExists(shortURL); ok {
		return ms.ShortToLongLinksMap[shortURL], nil
	}

	return "", nil
}

func (ms *MemoryStorage) StoreLink(link dto.Link) error {
	if _, ok := ms.ShortToLongLinksMap[link.Short]; ok {
		return &storageErrors.DuplicateError{}
	}

	if _, ok := ms.LongToShortLinksMap[link.Long]; ok {
		return &storageErrors.DuplicateError{}
	}

	err := ms.FileStorage.encoder.Encode(link)
	if err != nil {
		return err
	}

	ms.ShortToLongLinksMap[link.Short] = link.Long

	ms.LongToShortLinksMap[link.Long] = link.Short

	return nil
}

func (ms *MemoryStorage) StoreLinks(links []dto.Link) error {
	var ers error
	for _, link := range links {
		err := ms.StoreLink(link)
		if err != nil {
			ers = errors.Join(err)
		}
	}

	if ers != nil {
		return ers
	}

	return nil
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

		err = ms.StoreLink(*link)
		if err != nil {
			return err
		}
	}

	return nil
}
