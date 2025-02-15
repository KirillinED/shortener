package foundation

import (
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/services"
	"github.com/KirillinED/shortener/internal/storage"
	"github.com/KirillinED/shortener/internal/storage/interfaces"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Application interface {
	Bootstrap()
	Shutdown() error
	GetConfig() *config.Config
	GetLogger() *zap.Logger
	GetStorage() interfaces.Storage
	GetShortenerService() *services.ShortenerService
}

type App struct {
	Cfg              *config.Config
	Logger           *zap.Logger
	Storage          interfaces.Storage
	ShortenerService *services.ShortenerService
}

func NewApp() *App {
	app := &App{}

	app.Bootstrap()

	return app
}

func (app *App) Bootstrap() {
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

	app.Logger = zap.Must(zapConfig.Build())

	app.Cfg = config.NewConfig()

	s, err := storage.NewStorage(app.Cfg)
	if err != nil {
		app.Logger.Fatal(err.Error())
	}

	app.Storage = s

	app.ShortenerService = services.NewShortenerService(app.Cfg, s)
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

func (app *App) GetShortenerService() *services.ShortenerService {
	return app.ShortenerService
}
