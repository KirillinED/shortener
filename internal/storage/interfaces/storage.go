package interfaces

import "github.com/KirillinED/shortener/internal/dto"

type Storage interface {
	ShortExists(url string) (bool, error)
	LongExists(url string) (bool, error)
	GetShortURL(url string) (string, error)
	GetLongURL(url string) (string, error)
	StoreLink(link dto.Link) (bool, error)
	Close() error
}
