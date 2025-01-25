package storage

import (
	"errors"
	"github.com/KirillinED/shortener/internal/dto"
	"io"
)

type MemoryStorage struct {
	FileStorage         *FileStorage
	ShortToLongLinksMap map[string]string
	LongToShortLinksMap map[string]string
}

func NewMemoryStorage(fs *FileStorage) *MemoryStorage {
	return &MemoryStorage{
		FileStorage:         fs,
		ShortToLongLinksMap: make(map[string]string),
		LongToShortLinksMap: make(map[string]string),
	}
}

func (ms *MemoryStorage) ShortExists(shortURL string) bool {
	_, ok := ms.ShortToLongLinksMap[shortURL]

	return ok
}

func (ms *MemoryStorage) LongExists(longURL string) bool {
	_, ok := ms.LongToShortLinksMap[longURL]

	return ok
}

func (ms *MemoryStorage) GetShortURL(longURL string) string {
	if ms.LongExists(longURL) {
		return ms.LongToShortLinksMap[longURL]
	}

	return ""
}

func (ms *MemoryStorage) GetLongURL(shortURL string) string {
	if ms.ShortExists(shortURL) {
		return ms.ShortToLongLinksMap[shortURL]
	}

	return ""
}

func (ms *MemoryStorage) StoreLink(link dto.Link) error {
	err := ms.FileStorage.encoder.Encode(link)
	if err != nil {
		return err
	}

	ms.ShortToLongLinksMap[link.Short] = link.Long

	ms.LongToShortLinksMap[link.Long] = link.Short

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
