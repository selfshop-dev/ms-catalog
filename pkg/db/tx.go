package db

import (
	"context"
	"fmt"

	ctxval "github.com/selfshop-dev/lib-ctxval"
)

type TxFunc[T any] func(ctx context.Context, q T) error

type TxManager[T any] struct {
	Begin func(ctx context.Context) (
		T,
		func(ctx context.Context) error, // commit
		func(ctx context.Context) error, // rollback
		error,
	)
}

func (m *TxManager[T]) Do(ctx context.Context, fn TxFunc[T]) error {
	q, ok := ctxval.Get[T](ctx)
	if ok && any(q) != nil { // works for pointers, interfaces
		return fn(ctx, q)
	}

	q, com, rol, err := m.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}

	defer func() {
		_ = rol(ctx) //nolint:errcheck // rollback after error: best-effort, original error takes priority
	}()

	ctx = ctxval.With(ctx, q)

	if err := fn(ctx, q); err != nil {
		return err
	}

	if err := com(ctx); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
