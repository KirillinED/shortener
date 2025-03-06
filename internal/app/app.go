package app

import (
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/infrastructure"
	"github.com/KirillinED/shortener/internal/interfaces/links"
	"github.com/KirillinED/shortener/internal/interfaces/router"
	"github.com/KirillinED/shortener/internal/services/cookie"
	"github.com/KirillinED/shortener/internal/services/shortener"
	"github.com/KirillinED/shortener/internal/storage"
	"github.com/KirillinED/shortener/internal/storage/interfaces"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Application interface {
	Shutdown() error
	GetConfig() *config.Config
	GetLogger() *zap.Logger
	GetRouter() *chi.Mux
	GetLinkHandler() *links.LinkHandler
	// GetStorage() interfaces.Storage
	// GetShortenerService() *shortener.ShortenerService
	// GetSecureCookieService() *cookie.SecureCookie
}

type App struct {
	Cfg         *config.Config
	Logger      *zap.Logger
	Storage     interfaces.Storage
	Router      *router.Router
	LinkHandler *links.LinkHandler
	// ShortenerService    *shortener.ShortenerService
	// SecureCookieService *cookie.SecureCookie
}

func NewApp() *App {
	app := &App{}

	app.Cfg = config.NewConfig()

	s, err := storage.NewStorage(app.Cfg)
	if err != nil {
		app.Logger.Fatal(err.Error())
	}

	app.Storage = s

	app.Logger = infrastructure.NewLogger()

	NewStorage()
	NewLinkRepository(app.Storage)
	NewLinkService()
	app.LinkHandler = links.NewLinkHandler(app, NewLinkService)

	app.Router = router.NewRouter(app)
	app.Router.RegisterRoutes()

	// app.ShortenerService = shortener.NewShortenerService(app.Cfg, s)

	// app.SecureCookieService = cookie.NewSecureCookie(app.Cfg.CookieHashKey)

	return app
}

func (app *App) Shutdown() error {
	err := app.Logger.Sync()
	if err != nil {
		return err
	}

	err = app.Storage.Close()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) GetConfig() *config.Config {
	return app.Cfg
}

func (app *App) GetStorage() interfaces.Storage {
	return app.Storage
}

func (app *App) GetLogger() *zap.Logger {
	return app.Logger
}

func (app *App) GetRouter() *chi.Mux {
	return app.Router.GetRouter()
}

func (app *App) GetLinkHandler() *links.LinkHandler {
	return app.LinkHandler
}

func (app *App) GetShortenerService() *shortener.ShortenerService {
	return app.ShortenerService
}

func (app *App) GetSecureCookieService() *cookie.SecureCookie {
	return app.SecureCookieService
}
