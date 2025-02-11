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

		ok, err := app.GetStorage().ShortExists(requestBody.URL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		var shortUrl string
		if ok {
			shortUrl, err = app.GetStorage().GetShortURL(requestBody.URL)
		} else {
			shortUrl = utils.ShortURL(requestBody.URL)
			ok, err = app.GetStorage().StoreLink(dto.Link{
				Short: shortUrl,
				Long:  requestBody.URL,
			})

			if !ok {
				http.Error(w, "Something went wrong. Link is not save.", http.StatusInternalServerError)
			}
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		response := CreateShortLinkResponse{Result: app.GetConfig().BaseURL + shortUrl}

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

		ok, err := app.GetStorage().ShortExists(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		url, err := app.GetStorage().GetLongURL(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	}
}
