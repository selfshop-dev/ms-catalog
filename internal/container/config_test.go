package container_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	config "github.com/selfshop-dev/lib-config"

	"github.com/selfshop-dev/ms-catalog/internal/container"
)

func validDefaults() map[string]any {
	return map[string]any{
		"app.name":    "catalog",
		"app.runmode": "dev",

		"log.min_level": "info",
		"log.format":    "console",

		"entry.http.port":            8080,
		"entry.http.read_timeout":    "10s",
		"entry.http.request_timeout": "30s",
		"entry.http.write_timeout":   "35s",
		"entry.http.idle_timeout":    "60s",

		"db.dsn":                "postgres://localhost/catalog",
		"db.max_conns":          10,
		"db.min_conns":          2,
		"db.max_conn_lifetime":  "30m",
		"db.max_conn_idle_time": "5m",
	}
}

func TestConfig_valid(t *testing.T) {
	t.Parallel()
	_, err := config.New[container.Config]("APP", validDefaults())
	require.NoError(t, err)
}

func TestConfig_db(t *testing.T) {
	t.Parallel()

	t.Run("missing_dsn", func(t *testing.T) {
		t.Parallel()
		d := validDefaults()
		delete(d, "db.dsn")
		_, err := config.New[container.Config]("APP", d)
		require.Error(t, err)
		assert.ErrorContains(t, err, "Db.DSN")
	})

	t.Run("min_conns_exceeds_max_conns", func(t *testing.T) {
		t.Parallel()
		d := validDefaults()
		d["db.min_conns"] = 10
		d["db.max_conns"] = 5
		_, err := config.New[container.Config]("APP", d)
		require.Error(t, err)
		assert.ErrorContains(t, err, "Db.MinConns")
	})

	t.Run("idle_time_exceeds_lifetime", func(t *testing.T) {
		t.Parallel()
		d := validDefaults()
		d["db.max_conn_idle_time"] = "45m"
		d["db.max_conn_lifetime"] = "30m"
		_, err := config.New[container.Config]("APP", d)
		require.Error(t, err)
		assert.ErrorContains(t, err, "Db.MaxConnIdleTime")
	})

	t.Run("max_conns_out_of_range", func(t *testing.T) {
		t.Parallel()
		d := validDefaults()
		d["db.max_conns"] = 200
		_, err := config.New[container.Config]("APP", d)
		require.Error(t, err)
		assert.ErrorContains(t, err, "Db.MaxConns")
	})
}

func TestConfig_base_semantic(t *testing.T) {
	t.Parallel()

	t.Run("debug_forbidden_in_prod", func(t *testing.T) {
		t.Parallel()
		d := validDefaults()
		d["app.runmode"] = "prod"
		d["log.min_level"] = "info"
		d["debug"] = true
		_, err := config.New[container.Config]("APP", d)
		require.Error(t, err)
		assert.ErrorContains(t, err, "debug mode must be disabled in prod runmode")
	})
}
