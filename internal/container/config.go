package container

import (
	"time"

	config "github.com/selfshop-dev/lib-config"
)

type Config struct {
	Db          DbConfig `koanf:"db"`
	config.Base `koanf:",squash"`
}

type DbConfig struct {
	// DSN is the PostgreSQL connection string.
	DSN string `koanf:"dsn" validate:"required"`
	// MaxConns is the maximum number of connections in the pool.
	MaxConns int32 `koanf:"max_conns" validate:"required,min=1,max=100"`
	// MinConns is the minimum number of connections kept open.
	// Must be less than MaxConns.
	MinConns int32 `koanf:"min_conns" validate:"min=0,ltfield=MaxConns"`
	// MaxConnLifetime is the maximum duration a connection may be reused.
	MaxConnLifetime time.Duration `koanf:"max_conn_lifetime" validate:"required,gte=1m,lte=1h"`
	// MaxConnIdleTime is the maximum duration a connection may sit idle.
	// Must be less than MaxConnLifetime.
	MaxConnIdleTime time.Duration `koanf:"max_conn_idle_time" validate:"required,gte=30s,lte=30m,ltfield=MaxConnLifetime"`
}

func (c *Config) Validate() error {
	return c.Base.Validate()
}
