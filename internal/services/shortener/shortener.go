package shortener

import (
	"errors"
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/entities"
	storageErrors "github.com/KirillinED/shortener/internal/storage/errors"
	"github.com/KirillinED/shortener/internal/storage/interfaces"
	"github.com/KirillinED/shortener/internal/utils"
)

type ShortenerService struct {
	cfg     *config.Config
	storage interfaces.Storage
}

func NewShortenerService(cfg *config.Config, storage interfaces.Storage) *ShortenerService {
	return &ShortenerService{cfg: cfg, storage: storage}
}

func (s *ShortenerService) CreateShortLink(longUrl string) (string, error) {
	link := entities.Link{Long: longUrl, Short: utils.ShortURL(longUrl)}

	err := s.storage.StoreLink(link)
	if err != nil {
		return "", err
	}

	return s.cfg.BaseURL + link.Short, nil
}

func (s *ShortenerService) GetLongLink(shortUrl string) (string, error) {
	longUrl, err := s.storage.GetLongURL(shortUrl)
	if err != nil {
		return "", err
	}

	if longUrl == "" {
		return longUrl, &storageErrors.NotFoundError{}
	}

	return longUrl, nil
}

type CreateShortLinksResult struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (s *ShortenerService) CreateShortLinks(longUrls []struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}) ([]CreateShortLinksResult, error) {
	var links []entities.Link
	for _, item := range longUrls {
		links = append(links, entities.Link{
			CorrelationID: item.CorrelationID,
			Long:          item.OriginalURL,
			Short:         utils.ShortURL(item.OriginalURL),
		})
	}

	err := s.storage.StoreLinks(links)
	if err != nil && !errors.Is(err, &storageErrors.DuplicateError{}) {
		return nil, err
	}

	var res []CreateShortLinksResult
	for _, link := range links {
		res = append(res, CreateShortLinksResult{
			CorrelationID: link.CorrelationID,
			ShortURL:      s.cfg.BaseURL + link.Short,
		})
	}

	return res, err
}
