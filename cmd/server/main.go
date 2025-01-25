package main

import (
	"fmt"
	"github.com/KirillinED/shortener/internal/foundation"
	"github.com/KirillinED/shortener/internal/handlers"
	"github.com/KirillinED/shortener/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	app := foundation.NewApp()
	defer func() {
		err := app.Shutdown()
		if err != nil {
			log.Fatalf("failed to sync zap logger: %v", err)
		}
	}()

	if err := run(app); err != nil {
		panic(err)
	}
}

func run(app *foundation.App) error {
	r := chi.NewRouter()
	r.Use(
		middlewares.Logger(app.Logger),
		middlewares.Compress,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second))

	r.Post("/api/shorten", handlers.CreateShortLinkHandler(app))

	r.Get("/{link}", handlers.GetShortLinkHandler(app))

	fmt.Println("Listening on " + app.Cfg.Address.String())
	return http.ListenAndServe(app.Cfg.Address.String(), r)
}
