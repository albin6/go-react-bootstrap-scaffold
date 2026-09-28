package health_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"app/backend/internal/health"
	"app/backend/internal/http/response"
)

type mockPinger struct {
	err error
}

func (m *mockPinger) Ping(_ context.Context) error {
	return m.err
}

func TestHealthHandler_Liveness(t *testing.T) {
	svc := health.NewService(&mockPinger{})
	handler := health.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health/live", nil)
	rec := httptest.NewRecorder()

	handler.HandleLiveness(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res response.SuccessEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success to be true")
	}
}

func TestHealthHandler_Readiness_Healthy(t *testing.T) {
	svc := health.NewService(&mockPinger{err: nil})
	handler := health.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health/ready", nil)
	rec := httptest.NewRecorder()

	handler.HandleReadiness(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res response.SuccessEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success to be true")
	}
}

func TestHealthHandler_Readiness_Unhealthy(t *testing.T) {
	svc := health.NewService(&mockPinger{err: errors.New("connection refused")})
	handler := health.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health/ready", nil)
	rec := httptest.NewRecorder()

	handler.HandleReadiness(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status 503, got %d", rec.Code)
	}

	var res response.SuccessEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if res.Success {
		t.Errorf("expected success to be false on failed readiness")
	}
}

func TestHealthHandler_Aggregate(t *testing.T) {
	svc := health.NewService(&mockPinger{err: nil})
	handler := health.NewHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	rec := httptest.NewRecorder()

	handler.HandleAggregate(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var res response.SuccessEnvelope
	if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !res.Success {
		t.Errorf("expected success to be true")
	}
}
