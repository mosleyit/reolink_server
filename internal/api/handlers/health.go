package handlers

import (
	"net/http"
	"time"

	"github.com/mosleyit/reolink_server/pkg/utils"
)

var startTime = time.Now()

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HealthCheck handles health check requests
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	health := map[string]interface{}{
		"status":  "healthy",
		"version": "1.0.0",
		"uptime":  time.Since(startTime).String(),
	}

	utils.RespondJSON(w, http.StatusOK, health)
}

// ReadinessCheck handles readiness check requests
func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	// TODO: Check database, redis, and other dependencies
	ready := map[string]interface{}{
		"status": "ready",
		"components": map[string]string{
			"database": "healthy",
			"redis":    "healthy",
		},
	}

	utils.RespondJSON(w, http.StatusOK, ready)
}
