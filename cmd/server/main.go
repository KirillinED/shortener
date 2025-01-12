package main

import (
	"fmt"
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/handlers"
	"github.com/KirillinED/shortener/internal/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"net/http"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	zapConfig := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
			TimeKey:      "time",
			EncodeTime:   zapcore.RFC3339TimeEncoder,
		},
	}

	logger := zap.Must(zapConfig.Build())

	defer func() {
		err := logger.Sync()
		if err != nil {
			log.Fatalf("failed to sync zap logger: %v", err)
		}
	}()

	cfg := config.GetConfig()

	r := chi.NewRouter()
	r.Use(
		middlewares.Logger(logger),
		middlewares.Compress,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.Timeout(60*time.Second))

	r.Post("/api/shorten", handlers.CreateShortLinkHandler)

	r.Get("/{link}", handlers.GetShortLinkHandler)

	fmt.Println("Listening on " + cfg.Address.String())
	return http.ListenAndServe(cfg.Address.String(), r)
}
