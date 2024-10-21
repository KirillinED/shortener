package handlers

import (
	"github.com/KirillinED/shortener/internal/storage"
	"github.com/KirillinED/shortener/internal/utils"
	"io"
	"net/http"
)

const Domain = "http://localhost:8080/"

func CreateShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	url := utils.URL(body)

	if _, ok := storage.LongToShortLinksMap[url.String()]; !ok {
		shortURL := url.Short()
		storage.LongToShortLinksMap[url.String()] = shortURL
		storage.ShortToLongLinksMap[shortURL] = url.String()
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(Domain + storage.LongToShortLinksMap[url.String()]))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	short := r.URL.RequestURI()[1:]

	if _, ok := storage.ShortToLongLinksMap[short]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Location", storage.ShortToLongLinksMap[short])
	w.WriteHeader(http.StatusTemporaryRedirect)
}
