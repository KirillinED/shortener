package main

import (
	"fmt"
	"github.com/KirillinED/shortener/internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Post("/", handlers.CreateShortLinkHandler)

	r.Get("/{link}", handlers.GetShortLinkHandler)

	fmt.Println("Listening on port 8080")
	panic(http.ListenAndServe(":8080", r))
}
