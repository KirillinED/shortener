package interfaces

import "github.com/KirillinED/shortener/internal/dto"

type Storage interface {
	GetShortURL(url string) (string, error)
	GetLongURL(url string) (string, error)
	StoreLink(link dto.Link) error
	StoreLinks(links []dto.Link) error
	Close() error
}
