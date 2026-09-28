package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"app/backend/internal/http/response"
	"app/backend/internal/logger"
)

// Recoverer intercepts panics, logs the stack trace safely, and responds with a sanitized 500 error.
func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log := logger.FromContext(r.Context())
				stack := string(debug.Stack())

				log.Error("panic recovered in http handler",
					"error", fmt.Sprintf("%v", rec),
					"stack", stack,
				)

				response.InternalError(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
