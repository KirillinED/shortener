package router

import (
	"github.com/KirillinED/shortener/internal/app"
	"github.com/KirillinED/shortener/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

type Router struct {
	app    *app.App
	router *chi.Mux
}

func NewRouter(app app.Application) *Router {
	r := chi.NewRouter()
	r.Use(
		middlewares.Logger(app.GetLogger()),
		middlewares.CookieAuth(app),
		middlewares.Compress,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second),
	)
	return &Router{router: r}
}

func (r *Router) RegisterRoutes() {
	r.router.Get("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.router.Get("/api/user/links", r.app.GetLinkHandler().GetUserLinks)

	r.router.Post("/api/shorten/batch", r.app.GetLinkHandler().BatchCreateLinks)

	r.router.Post("/api/shorten", r.app.GetLinkHandler().CreateShortLink)

	r.router.Get("/{link}", r.app.GetLinkHandler().GetShortLink)
}

func (r *Router) GetRouter() *chi.Mux {
	return r.router
}
