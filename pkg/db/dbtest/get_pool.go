package dbtest

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	once     sync.Once
	err      error
	sharedDB *pgxpool.Pool
)

// MustGetPool returns a shared *pgxpool.Pool backed by a single postgres
// testcontainer instance started once per test binary execution.
//
// schema is the full DDL SQL applied after the container is ready.
// Pass your embed.FS content or a raw SQL string:
//
//	dbtest.MustGetPool(t, migrations.CurrentSchemaSQL)
func MustGetPool(t *testing.T, schema []byte) *pgxpool.Pool {
	t.Helper()
	once.Do(func() { sharedDB, err = initializePool(schema) })
	if err != nil {
		t.Fatalf("dbtest: init failed: %v", err)
	}
	return sharedDB
}

func initializePool(schema []byte) (*pgxpool.Pool, error) {
	ctx := context.Background()

	ctr, err := postgres.Run(ctx, "postgres:18-alpine",
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("start container: %w", err)
	}

	dsn, err := ctr.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("connection string: %w", err)
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	if _, err = pool.Exec(ctx, string(schema)); err != nil {
		return nil, fmt.Errorf("apply schema: %w", err)
	}
	return pool, nil
}
