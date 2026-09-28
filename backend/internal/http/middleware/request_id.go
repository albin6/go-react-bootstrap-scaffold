package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type requestIDKey struct{}

var requestIDCtxKey = requestIDKey{}

const HeaderXRequestID = "X-Request-ID"

// RequestID attaches an X-Request-ID header and stores the correlation ID in the request context.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(HeaderXRequestID)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		// Set response header
		w.Header().Set(HeaderXRequestID, reqID)

		// Attach to context
		ctx := context.WithValue(r.Context(), requestIDCtxKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext extracts the correlation request ID from context.
func RequestIDFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(requestIDCtxKey).(string); ok {
		return id
	}
	return ""
}
