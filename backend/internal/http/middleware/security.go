package middleware

import (
	"net/http"
)

// SecurityHeaders applies baseline security HTTP headers to all outgoing responses.
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strict Transport Security (HSTS) - 2 years
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")

		// Prevent MIME-type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")

		// Legacy XSS filter
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Referrer policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

		// Content Security Policy (strict defaults)
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none';")

		next.ServeHTTP(w, r)
	})
}
