package e2etest

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	config "github.com/selfshop-dev/lib-config"
	logger "github.com/selfshop-dev/lib-logger"

	"github.com/selfshop-dev/ms-catalog/internal/app"
	"github.com/selfshop-dev/ms-catalog/internal/db/gen"
	"github.com/selfshop-dev/ms-catalog/migrations"
	"github.com/selfshop-dev/ms-catalog/pkg/db/dbtest"
	"github.com/selfshop-dev/ms-catalog/pkg/health"
)

// NewExpect builds a real httptest.Server backed by the application router
// and returns an httpexpect.Expect client pointed at it.
//
// Each call creates a fresh server — isolation between test files is guaranteed
// by the httptest.Server lifecycle tied to t.Cleanup.
// The underlying postgres container and pool are shared across all callers
// within the same test binary to avoid the cost of spinning up multiple containers.
func NewExpect(t *testing.T) *httpexpect.Expect {
	t.Helper()

	pool := dbtest.MustGetPool(t, migrations.CurrentSchemaSQL)
	conf := &config.Base{
		App: config.App{
			Name:    "ms-catalog-e2e",
			Runmode: config.AppRunmodeDev,
		},
		Log: config.Log{
			MinLevel: config.LogMinLevelInfo,
			Format:   config.LogFormatConsole,
		},
		Entry: config.Entry{
			HTTP: config.HTTP{
				RequestTimeout: 30 * time.Second,
			},
		},
	}

	rout := app.NewRouter(
		gen.New(pool),
		newNopLogger(t),
		conf,
		newNopHealth(t),
	)

	serv := httptest.NewServer(rout)
	t.Cleanup(serv.Close)

	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  serv.URL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewDebugPrinter(t, true),
		},
	})
}

func newNopLogger(t *testing.T) *logger.Logger {
	t.Helper()
	cfg := logger.DefaultConfig()
	cfg.Development = true
	cfg.Sink = zapcore.AddSync(testWriter{t})
	l, err := logger.New(cfg)
	if err != nil {
		t.Fatalf("e2etest: build logger: %v", err)
	}
	return l
}

type testWriter struct{ t *testing.T }

func (w testWriter) Write(p []byte) (int, error) {
	w.t.Log(string(p))
	return len(p), nil
}

func (w testWriter) Sync() error { return nil }

// newNopHealth returns a health.Collector with no checkers for use in tests.
func newNopHealth(t *testing.T) *health.Collector {
	t.Helper()
	cfg := health.DefaultConfig()
	h, err := health.New(cfg, zap.NewNop())
	if err != nil {
		t.Fatalf("e2etest: build health: %v", err)
	}
	return h
}
