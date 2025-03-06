package middlewares

import (
	"encoding/hex"
	"errors"
	"github.com/KirillinED/shortener/internal/app"
	"github.com/KirillinED/shortener/internal/services/cookie"
	"github.com/KirillinED/shortener/internal/utils"
	"net/http"
)

// CookieAuth устанавливает сессионную куку с id пользователя, если она отсутствует у пользователя
func CookieAuth(app app.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := cookie.GetUserIdCookie(r)
			if errors.Is(err, http.ErrNoCookie) {
				var userId string
				userId, err = utils.GenerateUserId()
				if err != nil {
					http.Error(w, "middleware CookieAuth error: "+err.Error(), http.StatusInternalServerError)
					return
				}

				encode := app.GetSecureCookieService().EncodeString(userId)

				cookie.SetUserId(w, userId, hex.EncodeToString(encode))
			}

			if err != nil {
				http.Error(w, "middleware CookieAuth error: "+err.Error(), http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
