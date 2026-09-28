package routes

import (
	"net/http"

	"app/backend/internal/config"
	"app/backend/internal/health"
	"app/backend/internal/http/middleware"
)

// Setup builds and returns the application's root HTTP handler with all routes and middleware wired.
func Setup(cfg *config.Config, healthHandler *health.Handler) http.Handler {
	mux := http.NewServeMux()

	// Register Domain Module Routes
	healthHandler.RegisterRoutes(mux)

	// Apply Middlewares in execution order:
	// Outermost middleware executes first on inbound requests.
	var handler http.Handler = mux

	// 1. Enforce max request body size (1 MB default)
	handler = middleware.BodyLimit(1024 * 1024)(handler)

	// 2. Enforce CORS policies
	handler = middleware.CORS(cfg.CORSAllowedOrigins)(handler)

	// 3. Inject standard security headers
	handler = middleware.SecurityHeaders(handler)

	// 4. Structured HTTP logging
	handler = middleware.Logger(handler)

	// 5. Panic recovery
	handler = middleware.Recoverer(handler)

	// 6. Request correlation ID (outermost to ensure ID is available in all logging & recovery)
	handler = middleware.RequestID(handler)

	return handler
}
