package health

import (
	"net/http"

	"app/backend/internal/http/response"
)

// Handler handles HTTP requests for health check endpoints.
type Handler struct {
	service *Service
}

// NewHandler constructs a new health Handler.
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// RegisterRoutes mounts health endpoints onto the provided ServeMux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/health", h.HandleAggregate)
	mux.HandleFunc("GET /api/v1/health/live", h.HandleLiveness)
	mux.HandleFunc("GET /api/v1/health/ready", h.HandleReadiness)
}

// HandleLiveness returns a 200 OK if the application process is running.
func (h *Handler) HandleLiveness(w http.ResponseWriter, r *http.Request) {
	result := h.service.CheckLiveness(r.Context())
	response.Success(w, http.StatusOK, result)
}

// HandleReadiness returns 200 OK if downstream dependencies are ready, or 503 if not.
func (h *Handler) HandleReadiness(w http.ResponseWriter, r *http.Request) {
	result, ready := h.service.CheckReadiness(r.Context())
	if !ready {
		response.JSON(w, http.StatusServiceUnavailable, response.SuccessEnvelope{
			Success: false,
			Data:    result,
		})
		return
	}
	response.Success(w, http.StatusOK, result)
}

// HandleAggregate returns full system health metrics.
func (h *Handler) HandleAggregate(w http.ResponseWriter, r *http.Request) {
	result, healthy := h.service.CheckAggregate(r.Context())
	status := http.StatusOK
	if !healthy {
		status = http.StatusServiceUnavailable
	}
	response.JSON(w, status, response.SuccessEnvelope{
		Success: healthy,
		Data:    result,
	})
}
