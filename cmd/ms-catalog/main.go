// @title			selfshop-dev/ms-catalog API
// @version			0.1.0-dev
// @description	REST API для управления каталогом товаров и категорий.
//
// @contact.name	imidll
// @contact.url	https://github.com/selfshop-dev/ms-catalog
//
// @license.name	MIT
// @license.url	https://github.com/selfshop-dev/ms-catalog/blob/main/LICENSE
//
// @BasePath		/api/v1
//
// @schemes			http
package main

import (
	"fmt"

	_ "github.com/joho/godotenv/autoload"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	config "github.com/selfshop-dev/lib-config"
	logger "github.com/selfshop-dev/lib-logger"

	"github.com/selfshop-dev/ms-catalog/internal/app"
	"github.com/selfshop-dev/ms-catalog/internal/container"
	"github.com/selfshop-dev/ms-catalog/internal/db/gen"
	"github.com/selfshop-dev/ms-catalog/pkg/db"
	"github.com/selfshop-dev/ms-catalog/pkg/health"
	"github.com/selfshop-dev/ms-catalog/pkg/server"
)

func main() {
	fx.New(
		fx.Provide(newConfig,
			func(c *container.Config) *config.Base { return &c.Base },
			func(c *container.Config) *container.DbConfig { return &c.Db },
		),
		fx.Provide(newLogger), fx.WithLogger(func(l *logger.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: l.Unsampled().WithOptions(zap.IncreaseLevel(zap.WarnLevel), zap.AddCallerSkip(-1))}
		}),
		fx.Provide(newDb, func(db *db.Db) *gen.Queries { return gen.New(db) }),
		fx.Provide(
			fx.Annotate(newHealth, fx.ParamTags(`group:"checkers"`)),
			fx.Annotate(func(db *db.Db) health.Checker { return db }, fx.ResultTags(`group:"checkers"`)),
		),
		fx.Provide(app.NewRouter),
		fx.Provide(newServer), fx.Invoke(func(_ *server.Server) {}), // must be triggered
	).Run()
}

func newHealth(cs []health.Checker, l *logger.Logger) (*health.Collector, error) {
	d := health.DefaultConfig()
	d.MaxConcurrency = 10
	h, err := health.New(d, l.Unsampled(), cs...)
	if err != nil {
		return nil, fmt.Errorf("health: %w", err)
	}
	return h, nil
}

func newDb(lc fx.Lifecycle, c *container.DbConfig, l *logger.Logger) (*db.Db, error) {
	d := db.DefaultConfig()
	d.DSN = c.DSN
	d.MaxConns = c.MaxConns
	d.MinConns = c.MinConns
	d.MaxConnLifetime = c.MaxConnLifetime
	d.MaxConnIdleTime = c.MaxConnIdleTime
	e, err := db.New(d, l.Unsampled())
	if err != nil {
		return nil, fmt.Errorf("db: %w", err)
	}
	lc.Append(fx.StartStopHook(e.Start, e.Grace))
	return e, nil
}

func newConfig() (*container.Config, error) {
	c, err := config.New[container.Config]("INIT__", defaultValues)
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	return c, nil
}

func newServer(lc fx.Lifecycle, c *container.Config, l *logger.Logger, r *chi.Mux) (*server.Server, error) {
	d := server.DefaultConfig()
	d.Port = c.Entry.HTTP.Port
	d.ReadTimeout = c.Entry.HTTP.ReadTimeout
	d.IdleTimeout = c.Entry.HTTP.IdleTimeout
	d.WriteTimeout = c.Entry.HTTP.WriteTimeout
	s, err := server.NewServer(d, l, r)
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}
	lc.Append(fx.StartStopHook(s.Start, s.Grace))
	return s, nil
}

func newLogger(lc fx.Lifecycle, c *config.Base) (*logger.Logger, error) {
	p, err := zapcore.ParseLevel(c.Log.MinLevel)
	if err != nil {
		return nil, fmt.Errorf("logger: log level parse: %w", err)
	}
	d := logger.DefaultConfig()
	d.Version = version
	d.InitialFields = map[string]string{
		"comhash": comhash,
		"buildAt": buildAt,
	}
	d.Development = c.IsDev()
	d.Level = logger.NewLevelManager(p)
	l, err := logger.New(d)
	if err != nil {
		return nil, fmt.Errorf("logger: %w", err)
	}
	lc.Append(fx.StopHook(l.Sync))
	return l, nil
}
