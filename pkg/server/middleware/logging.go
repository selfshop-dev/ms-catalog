package middleware

import (
	"net"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func Logging(l *zap.Logger) func(http.Handler) http.Handler {
	l = l.Named("server.access")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			ip := r.RemoteAddr
			if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
				ip = host
			}

			fields := []zap.Field{
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.String("query", r.URL.RawQuery),
				zap.Int("status", rw.statusCode),
				zap.Int("bytes", rw.bytesWritten),
				zap.Duration("latency", duration),
				zap.String("remote_ip", ip),
				zap.String("user_agent", r.UserAgent()),
				zap.String("request_id", RequestIDFromContext(r.Context())),
			}

			switch {
			case rw.statusCode >= 500:
				l.Error("server error", fields...)
			case rw.statusCode >= 400:
				l.Warn("client error", fields...)
			default:
				l.Info("request handled", fields...)
			}
		})
	}
}

type responseWriter struct {
	http.ResponseWriter

	statusCode   int
	bytesWritten int
	written      bool
}

func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *responseWriter) Write(bs []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	n, err := rw.ResponseWriter.Write(bs)
	rw.bytesWritten += n
	return n, err
}
