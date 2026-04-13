package db

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

func (db *Db) Start(ctx context.Context) error {
	const pingTimeout = 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	stat := db.Stat()
	db.l.Info("connected",
		zap.Int32("total_conns", stat.TotalConns()),
		zap.Int32("max_conns", stat.MaxConns()),
	)
	return nil
}

func (db *Db) Grace(_ context.Context) error {
	db.Close()
	db.l.Info("pool closed")
	return nil
}
