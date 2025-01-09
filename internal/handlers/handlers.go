package handlers

import (
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/storage"
	"github.com/KirillinED/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

type Map struct {
	arr []any
}

func CreateShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "body cannot be empty", http.StatusBadRequest)
		return
	}

	url := utils.URL(body)

	if _, ok := storage.LongToShortLinksMap[url.String()]; !ok {
		shortURL := url.Short()
		storage.LongToShortLinksMap[url.String()] = shortURL
		storage.ShortToLongLinksMap[shortURL] = url.String()
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(config.GetConfig().BaseURL + storage.LongToShortLinksMap[url.String()]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	link := chi.URLParam(r, "link")

	if _, ok := storage.ShortToLongLinksMap[link]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", storage.ShortToLongLinksMap[link])
	w.WriteHeader(http.StatusTemporaryRedirect)
}
