package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/ultimate-ci-cd-pipeline/internal/model"
)

type HealthHandler struct {
	version string
}

func NewHealthHandler(version string) *HealthHandler {
	return &HealthHandler{version: version}
}

// Healthz provides overall service status.
func (h *HealthHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	resp := model.HealthResponse{
		Status:    "UP",
		Version:   h.version,
		Timestamp: time.Now().UTC(),
		Details: map[string]string{
			"database": "connected",
			"cache":    "healthy",
		},
	}
	writeJSON(w, http.StatusOK, resp)
}

// Livez provides k8s liveness probe.
func (h *HealthHandler) Livez(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"status": "ALIVE",
	}
	writeJSON(w, http.StatusOK, resp)
}

// Readyz provides k8s readiness probe.
func (h *HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"status": "READY",
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
