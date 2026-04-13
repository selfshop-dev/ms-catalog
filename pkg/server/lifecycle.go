package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go s.serve(errCh)

	const startTimeout = 100 * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, startTimeout)
	defer cancel()

	select {
	case <-ctx.Done():
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			s.l.Info("started")
			return nil
		}
		return fmt.Errorf("start: %w", ctx.Err())

	case err := <-errCh:
		return fmt.Errorf("start: %w", err)
	}
}

func (s *Server) serve(errCh chan<- error) {
	if err := s.s.Serve(s.ln); !errors.Is(err, http.ErrServerClosed) {
		errCh <- err
	}
}

func (s *Server) Grace(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, s.shutdownTimeout)
	defer cancel()

	s.l.Info("shutdown started", zap.Duration("timeout", s.shutdownTimeout))

	if err := s.s.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	s.l.Info("stopped")
	return nil
}
