package health

import "time"

// LivenessResponse represents the payload for process liveness probes.
type LivenessResponse struct {
	Status string `json:"status"`
}

// ReadinessResponse represents the payload for dependency readiness checks.
type ReadinessResponse struct {
	Status   string            `json:"status"`
	Database string            `json:"database"`
	Details  map[string]string `json:"details,omitempty"`
}

// HealthResponse represents the aggregate health check payload.
type HealthResponse struct {
	Status    string            `json:"status"`
	Uptime    string            `json:"uptime"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}
