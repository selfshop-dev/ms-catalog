package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Db struct {
	*pgxpool.Pool
	l *zap.Logger
}

func New(c Config, l *zap.Logger) (*Db, error) {
	d, err := pgxpool.ParseConfig(c.DSN)
	if err != nil {
		return nil, fmt.Errorf("parse dsn: %w", err)
	}

	d.MaxConns = c.MaxConns
	d.MinConns = c.MinConns
	d.MaxConnLifetime = c.MaxConnLifetime
	d.MaxConnIdleTime = c.MaxConnIdleTime

	l = l.Named("db::postgres")

	p, err := pgxpool.NewWithConfig(context.Background(), d)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	l.Info("initialized")

	return &Db{
		Pool: p,
		l:    l,
	}, nil
}
