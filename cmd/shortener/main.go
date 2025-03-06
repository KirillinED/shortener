package main

import (
	"fmt"
	"github.com/KirillinED/shortener/internal/app"
	"log"
	"net/http"
)

func main() {
	application := app.NewApp()
	defer func() {
		err := application.Shutdown()
		if err != nil {
			log.Fatalf("failed to sync zap logger: %v", err)
		}
	}()

	fmt.Println("Listening on " + application.Cfg.Address.String())
	if err := http.ListenAndServe(application.Cfg.Address.String(), application.GetRouter()); err != nil {
		panic(err)
	}
}
