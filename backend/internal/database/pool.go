package database

import (
	"context"
	"fmt"
	"time"

	"app/backend/internal/config"
	"app/backend/internal/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool creates and initializes a pgxpool.Pool configured according to application config.
func NewPool(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	log := logger.FromContext(ctx)

	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database url config: %w", err)
	}

	// Apply connection pool sizing and lifecycle settings
	poolConfig.MaxConns = int32(cfg.DBMaxOpenConns)
	poolConfig.MinConns = int32(cfg.DBMinOpenConns)
	poolConfig.MaxConnIdleTime = cfg.DBMaxIdleTime
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	log.Info("initializing postgresql connection pool",
		"max_conns", poolConfig.MaxConns,
		"min_conns", poolConfig.MinConns,
		"max_idle_time", poolConfig.MaxConnIdleTime,
	)

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	// Ping database with a reasonable connection verification timeout
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		// In development, we don't necessarily want hard crash if local docker isn't started yet,
		// but we log a clear warning. The health readiness probe will report degraded.
		log.Warn("database ping failed during startup; database may be offline or starting up",
			"error", err.Error(),
		)
	} else {
		log.Info("successfully connected and verified postgresql database pool")
	}

	return pool, nil
}
