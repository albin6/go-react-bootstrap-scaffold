package health

import (
	"context"
	"time"
)

// DatabasePinger is a minimal consumer-defined interface for checking database connectivity.
type DatabasePinger interface {
	Ping(ctx context.Context) error
}

// Service evaluates application health and dependency statuses.
type Service struct {
	db        DatabasePinger
	startTime time.Time
}

// NewService constructs a new HealthService.
func NewService(db DatabasePinger) *Service {
	return &Service{
		db:        db,
		startTime: time.Now(),
	}
}

// CheckLiveness validates that the Go runtime process is running and responding.
func (s *Service) CheckLiveness(_ context.Context) LivenessResponse {
	return LivenessResponse{
		Status: "alive",
	}
}

// CheckReadiness validates that downstream dependencies (e.g. PostgreSQL) are accessible.
func (s *Service) CheckReadiness(ctx context.Context) (ReadinessResponse, bool) {
	if s.db == nil {
		return ReadinessResponse{
			Status:   "not_ready",
			Database: "unconfigured",
		}, false
	}

	pingCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := s.db.Ping(pingCtx); err != nil {
		return ReadinessResponse{
			Status:   "not_ready",
			Database: "down",
			Details: map[string]string{
				"error": err.Error(),
			},
		}, false
	}

	return ReadinessResponse{
		Status:   "ready",
		Database: "up",
	}, true
}

// CheckAggregate returns the overall system health, including uptime and all sub-checks.
func (s *Service) CheckAggregate(ctx context.Context) (HealthResponse, bool) {
	readiness, ready := s.CheckReadiness(ctx)

	overallStatus := "healthy"
	if !ready {
		overallStatus = "degraded"
	}

	return HealthResponse{
		Status:    overallStatus,
		Uptime:    time.Since(s.startTime).Truncate(time.Second).String(),
		Timestamp: time.Now().UTC(),
		Checks: map[string]string{
			"process":  "alive",
			"database": readiness.Database,
		},
	}, ready
}
