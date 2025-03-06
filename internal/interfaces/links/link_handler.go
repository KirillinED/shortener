package links

import (
	"encoding/json"
	"errors"
	"github.com/KirillinED/shortener/internal/app"
	"github.com/KirillinED/shortener/internal/services/cookie"
	storageErrors "github.com/KirillinED/shortener/internal/storage/errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type LinkHandler struct {
	app app.Application
}

func NewLinkHandler(app app.Application) *LinkHandler {
	return &LinkHandler{app: app}
}

func (h *LinkHandler) GetUserLinks(w http.ResponseWriter, r *http.Request) {
	var statusCode = http.StatusNoContent

	u, err := cookie.GetUserIdCookie(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ok, err := h.app.GetSecureCookieService().IsValidUserIdCookie(u)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		statusCode = http.StatusUnauthorized
	}

	// @TODO get users links
	links := h.app.GetStorage().GetUserLinks(u.UserId)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
	w.WriteHeader(statusCode)
}

func (h *LinkHandler) BatchCreateLinks(w http.ResponseWriter, r *http.Request) {
	req := make([]struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}, 0)

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	statusCode := http.StatusCreated
	links, err := h.app.GetShortenerService().CreateShortLinks(req)
	if err != nil {
		if !errors.Is(err, &storageErrors.DuplicateError{}) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		statusCode = http.StatusConflict
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(links)
	if err != nil {
		h.app.GetLogger().Error("response encode error: " + err.Error())
	}
}

func (h *LinkHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shortUrl, err := h.app.GetShortenerService().CreateShortLink(req.URL)
	statusCode := http.StatusCreated
	if err != nil {
		if errors.Is(err, &storageErrors.DuplicateError{}) {
			statusCode = http.StatusConflict
		} else {
			http.Error(w, "Something went wrong. Link is not save.", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err = json.NewEncoder(w).Encode(map[string]string{"result": shortUrl})
	if err != nil {
		h.app.GetLogger().Error("response encode error: " + err.Error())
	}
}

func (h *LinkHandler) GetShortLink(w http.ResponseWriter, r *http.Request) {
	url, err := app.GetShortenerService().GetLongLink(chi.URLParam(r, "link"))
	if err != nil {
		if errors.Is(err, &storageErrors.NotFoundError{}) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
