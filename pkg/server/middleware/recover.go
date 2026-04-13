package middleware

import (
	"net/http"
	"runtime/debug"

	"go.uber.org/zap"
)

func Recover(l *zap.Logger) func(http.Handler) http.Handler {
	l = l.Named("server.recovery")
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					l.Error("panic recovered",
						zap.Any("panic", rec),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.String("request_id", RequestIDFromContext(r.Context())),
						zap.ByteString("stack", debug.Stack()),
					)
					http.Error(w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
