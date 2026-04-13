package db

import (
	"context"

	"github.com/selfshop-dev/ms-catalog/pkg/health"
)

var _ health.Checker = (*Db)(nil)

func (db *Db) Name() string                    { return "db::postgres" }
func (db *Db) Check(ctx context.Context) error { return db.Ping(ctx) }
