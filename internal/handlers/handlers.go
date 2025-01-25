package handlers

import (
	"encoding/json"
	"github.com/KirillinED/shortener/internal/dto"
	"github.com/KirillinED/shortener/internal/foundation"
	"github.com/KirillinED/shortener/internal/utils"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
)

func CreateShortLinkHandler(app foundation.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		if !app.GetMemoryStorage().ShortExists(requestBody.URL) {
			err = app.GetMemoryStorage().StoreLink(dto.Link{
				Short: utils.ShortURL(requestBody.URL),
				Long:  requestBody.URL,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		response := CreateShortLinkResponse{Result: app.GetConfig().BaseURL + app.GetMemoryStorage().GetShortURL(requestBody.URL)}

		res, err := json.Marshal(response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = w.Write(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func GetShortLinkHandler(app foundation.Application) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		link := chi.URLParam(r, "link")

		if !app.GetMemoryStorage().ShortExists(link) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Location", app.GetMemoryStorage().GetLongURL(link))
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
