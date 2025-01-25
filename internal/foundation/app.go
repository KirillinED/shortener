package foundation

import (
	"github.com/KirillinED/shortener/internal/config"
	"github.com/KirillinED/shortener/internal/storage"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Application interface {
	Bootstrap()
	Shutdown() error
	GetConfig() *config.Config
	GetLogger() *zap.Logger
	GetMemoryStorage() *storage.MemoryStorage
}

type App struct {
	Cfg           *config.Config
	Logger        *zap.Logger
	MemoryStorage *storage.MemoryStorage
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

	app.MemoryStorage = storage.NewMemoryStorage(storage.NewFileStorage(app.Cfg.FileStoragePath))

	err := app.MemoryStorage.Recovering()
	if err != nil {
		app.Logger.Fatal(err.Error())
	}
}

func (app *App) Shutdown() error {
	err := app.Logger.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) GetConfig() *config.Config {
	return app.Cfg
}

func (app *App) GetMemoryStorage() *storage.MemoryStorage {
	return app.MemoryStorage
}

func (app *App) GetLogger() *zap.Logger {
	return app.Logger
}

// AppStub struct for testing
type AppStub struct {
	Cfg           *config.Config
	Logger        *zap.Logger
	MemoryStorage *storage.MemoryStorage
}

func NewAppStub() *AppStub {
	app := &AppStub{}

	app.Bootstrap()

	return app
}

func (app *AppStub) Bootstrap() {
	app.Logger = zap.NewNop()

	app.Cfg = config.DefaultConfig()

	app.MemoryStorage = storage.NewMemoryStorage(storage.NewFileStorage(app.Cfg.FileStoragePath))
}

func (app *AppStub) Shutdown() error {
	err := app.Logger.Sync()
	if err != nil {
		return err
	}

	err = app.MemoryStorage.Close()
	if err != nil {
		return err
	}

	return nil
}

func (app *AppStub) GetConfig() *config.Config {
	return app.Cfg
}

func (app *AppStub) GetLogger() *zap.Logger {
	return app.Logger
}

func (app *AppStub) GetMemoryStorage() *storage.MemoryStorage {
	return app.MemoryStorage
}
