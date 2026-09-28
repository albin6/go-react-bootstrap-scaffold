package middleware

import (
	"net/http"
	"time"

	"app/backend/internal/logger"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// Logger creates a structured request logger middleware using log/slog.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := RequestIDFromContext(r.Context())

		// Create contextual logger enriched with request attributes
		ctxLogger := logger.FromContext(r.Context()).With(
			"request_id", reqID,
			"method", r.Method,
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr,
		)

		ctx := logger.WithContext(r.Context(), ctxLogger)

		wrapper := &responseWriterWrapper{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapper, r.WithContext(ctx))

		duration := time.Since(start)

		if wrapper.statusCode >= 500 {
			ctxLogger.Error("http request error",
				"status", wrapper.statusCode,
				"duration_ms", duration.Milliseconds(),
				"bytes", wrapper.written,
			)
		} else {
			ctxLogger.Info("http request completed",
				"status", wrapper.statusCode,
				"duration_ms", duration.Milliseconds(),
				"bytes", wrapper.written,
			)
		}
	})
}
