package middlewares

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}

			defer func() {
				// Логируем запрос
				logger.Info("Request",
					zap.String("method", r.Method),
					zap.String("uri", r.RequestURI),
					zap.Int("status", ww.status),
					zap.Int("size", ww.size),
					zap.Float64("duration", time.Since(start).Seconds()),
					zap.String("remote_addr", r.RemoteAddr),
					zap.String("user_agent", r.UserAgent()),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
