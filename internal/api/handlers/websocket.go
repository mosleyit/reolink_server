package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/mosleyit/reolink_server/internal/api/service"
	"github.com/mosleyit/reolink_server/internal/logger"
	"github.com/mosleyit/reolink_server/pkg/utils"
)

// EventStreamServiceInterface defines the interface for event streaming
type EventStreamServiceInterface interface {
	Subscribe(ctx context.Context, subscriberID string, cameraID string, bufferSize int) *service.EventStreamSubscriber
	Unsubscribe(subscriberID string)
}

// EventStreamHandler handles WebSocket and SSE connections for events
type EventStreamHandler struct {
	streamService EventStreamServiceInterface
}

// NewEventStreamHandler creates a new event stream handler
func NewEventStreamHandler(streamService EventStreamServiceInterface) *EventStreamHandler {
	return &EventStreamHandler{
		streamService: streamService,
	}
}

// WebSocketEvents handles WebSocket connections for all events
func (h *EventStreamHandler) WebSocketEvents(w http.ResponseWriter, r *http.Request) {
	h.handleWebSocket(w, r, "")
}

// WebSocketCameraEvents handles WebSocket connections for camera-specific events
func (h *EventStreamHandler) WebSocketCameraEvents(w http.ResponseWriter, r *http.Request) {
	cameraID := chi.URLParam(r, "id")
	if cameraID == "" {
		utils.RespondBadRequest(w, "Camera ID is required", nil)
		return
	}
	h.handleWebSocket(w, r, cameraID)
}

// handleWebSocket handles WebSocket connections
func (h *EventStreamHandler) handleWebSocket(w http.ResponseWriter, r *http.Request, cameraID string) {
	ctx := r.Context()

	// Upgrade to WebSocket
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // Allow connections from any origin in development
	})
	if err != nil {
		logger.Error("Failed to upgrade to WebSocket", zap.Error(err))
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "Connection closed")

	// Create subscriber
	subscriberID := uuid.New().String()
	subscriber := h.streamService.Subscribe(ctx, subscriberID, cameraID, 100)
	defer h.streamService.Unsubscribe(subscriberID)

	logger.Info("WebSocket client connected",
		zap.String("subscriber_id", subscriberID),
		zap.String("camera_id", cameraID))

	// Send events to client
	for {
		select {
		case event, ok := <-subscriber.EventCh:
			if !ok {
				return
			}

			// Marshal event to JSON
			data, err := json.Marshal(event)
			if err != nil {
				logger.Error("Failed to marshal event", zap.Error(err))
				continue
			}

			// Send to client
			err = conn.Write(ctx, websocket.MessageText, data)
			if err != nil {
				logger.Error("Failed to send WebSocket message", zap.Error(err))
				return
			}

		case <-ctx.Done():
			logger.Info("WebSocket context cancelled", zap.String("subscriber_id", subscriberID))
			return
		}
	}
}

// SSEEvents handles Server-Sent Events for all events
func (h *EventStreamHandler) SSEEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create subscriber
	subscriberID := uuid.New().String()
	subscriber := h.streamService.Subscribe(ctx, subscriberID, "", 100)
	defer h.streamService.Unsubscribe(subscriberID)

	logger.Info("SSE client connected", zap.String("subscriber_id", subscriberID))

	// Create flusher
	flusher, ok := w.(http.Flusher)
	if !ok {
		logger.Error("Streaming not supported")
		utils.RespondError(w, http.StatusInternalServerError, "STREAMING_NOT_SUPPORTED", "Streaming not supported", nil)
		return
	}

	// Send initial connection message
	fmt.Fprintf(w, "data: {\"type\":\"connected\",\"subscriber_id\":\"%s\"}\n\n", subscriberID)
	flusher.Flush()

	// Send heartbeat ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	// Send events to client
	for {
		select {
		case event, ok := <-subscriber.EventCh:
			if !ok {
				return
			}

			// Marshal event to JSON
			data, err := json.Marshal(event)
			if err != nil {
				logger.Error("Failed to marshal event", zap.Error(err))
				continue
			}

			// Send SSE event
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()

		case <-ticker.C:
			// Send heartbeat
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()

		case <-ctx.Done():
			logger.Info("SSE context cancelled", zap.String("subscriber_id", subscriberID))
			return
		}
	}
}

// Legacy standalone functions for backward compatibility
// These will be replaced by the handler methods in the router

// WebSocketEvents handles WebSocket connections for all events
func WebSocketEvents(w http.ResponseWriter, r *http.Request) {
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "WebSocket events not yet implemented - use handler instance", nil)
}

// WebSocketCameraEvents handles WebSocket connections for camera-specific events
func WebSocketCameraEvents(w http.ResponseWriter, r *http.Request) {
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "WebSocket camera events not yet implemented - use handler instance", nil)
}

// SSEEvents handles Server-Sent Events for all events
func SSEEvents(w http.ResponseWriter, r *http.Request) {
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "SSE events not yet implemented - use handler instance", nil)
}
