package handlers

import (
	"encoding/json"
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/storage"
	"github.com/KirillinED/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func CreateShortLinkHandler(w http.ResponseWriter, r *http.Request) {
	type CreateShortLinkRequestBody struct {
		URL string `json:"url"`
	}

	type CreateShortLinkResponse struct {
		Result string `json:"result"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(body) == 0 {
		http.Error(w, "body cannot be empty", http.StatusBadRequest)
		return
	}

	requestBody := CreateShortLinkRequestBody{}
	if err = json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if _, ok := storage.LongToShortLinksMap[requestBody.URL]; !ok {
		shortURL := utils.ShortURL(requestBody.URL)
		storage.LongToShortLinksMap[requestBody.URL] = shortURL
		storage.ShortToLongLinksMap[shortURL] = requestBody.URL
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	response := CreateShortLinkResponse{Result: config.GetConfig().BaseURL + storage.LongToShortLinksMap[requestBody.URL]}

	res, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = w.Write(res)
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
