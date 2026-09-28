package response

import (
	"encoding/json"
	"net/http"
)

// SuccessEnvelope defines the standardized format for successful HTTP responses.
type SuccessEnvelope struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
}

// ErrorDetail represents specific validation or field-level issues.
type ErrorDetail struct {
	Field string `json:"field,omitempty"`
	Issue string `json:"issue"`
}

// ErrorInfo encapsulates standardized error payload metadata.
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details []any  `json:"details,omitempty"`
}

// ErrorEnvelope defines the standardized format for HTTP error responses.
type ErrorEnvelope struct {
	Success bool      `json:"success"`
	Error   ErrorInfo `json:"error"`
}

// JSON serializes data and writes it to the response writer with the specified status code.
func JSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

// Success renders a standardized 2xx success envelope.
func Success(w http.ResponseWriter, statusCode int, data any) {
	JSON(w, statusCode, SuccessEnvelope{
		Success: true,
		Data:    data,
	})
}

// Error renders a standardized 4xx/5xx error envelope.
func Error(w http.ResponseWriter, statusCode int, code, message string, details ...any) {
	JSON(w, statusCode, ErrorEnvelope{
		Success: false,
		Error: ErrorInfo{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// BadRequest writes a 400 Bad Request error response.
func BadRequest(w http.ResponseWriter, message string, details ...any) {
	Error(w, http.StatusBadRequest, "BAD_REQUEST", message, details...)
}

// NotFound writes a 404 Not Found error response.
func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, "NOT_FOUND", message)
}

// Conflict writes a 409 Conflict error response.
func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, "CONFLICT", message)
}

// InternalError writes a sanitized 500 Internal Server Error response.
func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", "An internal server error occurred.")
}

// ServiceUnavailable writes a 503 Service Unavailable response.
func ServiceUnavailable(w http.ResponseWriter, message string) {
	Error(w, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", message)
}
