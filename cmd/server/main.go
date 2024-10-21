package main

import (
	"fmt"
	"github.com/KirillinED/shortener/internal/handlers"
	"github.com/KirillinED/shortener/internal/router"
	"net/http"
)

func main() {
	r := router.NewRouter()
	r.Handle(http.MethodPost, "/", handlers.CreateShortLinkHandler)

	r.Handle(http.MethodGet, `/.*`, handlers.GetShortLinkHandler)

	fmt.Println("Listening on port 8080")
	panic(http.ListenAndServe(":8080", r))
}
