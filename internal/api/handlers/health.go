package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/mosleyit/reolink_server/pkg/utils"
)

var startTime = time.Now()

// HealthChecker interface for components that can report health
type HealthChecker interface {
	Ping() error
}

// HealthHandler handles health check requests
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
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
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	components := make(map[string]string)
	allHealthy := true

	// Check database
	if h.db != nil {
		if err := h.db.PingContext(ctx); err != nil {
			components["database"] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			components["database"] = "healthy"
		}
	} else {
		components["database"] = "not configured"
	}

	status := "ready"
	statusCode := http.StatusOK
	if !allHealthy {
		status = "degraded"
		statusCode = http.StatusServiceUnavailable
	}

	ready := map[string]interface{}{
		"status":     status,
		"components": components,
	}

	utils.RespondJSON(w, statusCode, ready)
}
