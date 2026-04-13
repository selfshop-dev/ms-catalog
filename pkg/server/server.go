package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	logger "github.com/selfshop-dev/lib-logger"
	"go.uber.org/zap"
)

type Server struct {
	s               *http.Server
	l               *logger.Logger
	ln              net.Listener
	shutdownTimeout time.Duration
}

func NewServer(c Config, l *logger.Logger, r http.Handler) (*Server, error) {
	lc := &net.ListenConfig{}
	t, err := lc.Listen(context.Background(), "tcp",
		net.JoinHostPort(c.Host, strconv.FormatUint(uint64(c.Port), 10)))
	if err != nil {
		return nil, fmt.Errorf("listen: %w", err)
	}

	l = l.Named("http::server").With(zap.String("addr", t.Addr().String()))

	s := &http.Server{
		Addr:              t.Addr().String(),
		Handler:           r,
		ReadTimeout:       c.ReadTimeout,
		ReadHeaderTimeout: c.ReadHeaderTimeout,
		IdleTimeout:       c.IdleTimeout,
		WriteTimeout:      c.WriteTimeout,
	}

	l.Info("initialized")

	return &Server{
		s:               s,
		ln:              t,
		l:               l,
		shutdownTimeout: c.ShutdownTimeout,
	}, nil
}
