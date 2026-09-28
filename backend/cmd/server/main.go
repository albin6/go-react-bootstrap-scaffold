package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app/backend/internal/config"
	"app/backend/internal/database"
	"app/backend/internal/health"
	"app/backend/internal/http/routes"
	"app/backend/internal/logger"
)

func main() {
	// 1. Load and validate environment configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize structured logging
	log := logger.Init(cfg.LogLevel, cfg.LogFormat)
	log.Info("starting backend service",
		"env", cfg.AppEnv,
		"port", cfg.HTTPPort,
		"log_level", cfg.LogLevel,
	)

	// Root context for initial bootstrapping
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 3. Initialize PostgreSQL database connection pool
	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		log.Error("failed to initialize database pool", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	// 4. Execute database migrations
	if err := database.Migrate(ctx, pool); err != nil {
		log.Error("database migration error", "error", err)
	}

	// 5. Construct domain services and HTTP handlers
	healthService := health.NewService(pool)
	healthHandler := health.NewHandler(healthService)

	// 6. Assemble HTTP router and middleware
	httpHandler := routes.Setup(cfg, healthHandler)

	// 7. Explicitly configure HTTP server timeouts
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler:           httpHandler,
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}

	// 8. Start HTTP server asynchronously
	go func() {
		log.Info("http server listening", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("http server stopped unexpectedly", "error", err)
			os.Exit(1)
		}
	}()

	// 9. Listen for OS termination signals (SIGINT, SIGTERM)
	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()
	log.Info("shutdown signal received; commencing graceful teardown")

	// 10. Execute graceful shutdown with timeout
	graceCtx, graceCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer graceCancel()

	if err := server.Shutdown(graceCtx); err != nil {
		log.Error("graceful server shutdown failed", "error", err)
	} else {
		log.Info("http server stopped cleanly")
	}

	// Close database connection pool
	pool.Close()
	log.Info("database pool connections closed")

	log.Info("service terminated successfully")
}
