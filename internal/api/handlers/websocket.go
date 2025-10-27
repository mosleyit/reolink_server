package handlers

import (
	"net/http"

	"github.com/mosleyit/reolink_server/pkg/utils"
)

// WebSocketEvents handles WebSocket connections for all events
func WebSocketEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement WebSocket event streaming
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "WebSocket events not yet implemented", nil)
}

// WebSocketCameraEvents handles WebSocket connections for camera-specific events
func WebSocketCameraEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement WebSocket camera event streaming
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "WebSocket camera events not yet implemented", nil)
}

// SSEEvents handles Server-Sent Events for all events
func SSEEvents(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement SSE event streaming
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "SSE events not yet implemented", nil)
}

