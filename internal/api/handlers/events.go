package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/mosleyit/reolink_server/pkg/utils"
)

// EventServiceInterface defines the interface for event service operations
type EventServiceInterface interface {
	GetEvent(ctx context.Context, id string) (*models.Event, error)
	ListEvents(ctx context.Context, limit, offset int) ([]*models.Event, error)
	CountEvents(ctx context.Context) (int, error)
	AcknowledgeEvent(ctx context.Context, id string) error
}

// EventHandler handles event-related HTTP requests
type EventHandler struct {
	eventService  EventServiceInterface
	cameraService CameraServiceInterface
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventService EventServiceInterface, cameraService CameraServiceInterface) *EventHandler {
	return &EventHandler{
		eventService:  eventService,
		cameraService: cameraService,
	}
}

// ListEvents handles GET /api/v1/events
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit <= 0 {
		limit = 50
	}

	events, err := h.eventService.ListEvents(ctx, limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to list events", err)
		return
	}

	total, err := h.eventService.CountEvents(ctx)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to count events", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"events": events,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetEvent handles GET /api/v1/events/{id}
func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	event, err := h.eventService.GetEvent(ctx, id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "NOT_FOUND", "Event not found", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, event)
}

// AcknowledgeEvent handles PUT /api/v1/events/{id}/acknowledge
func (h *EventHandler) AcknowledgeEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if err := h.eventService.AcknowledgeEvent(ctx, id); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to acknowledge event", err)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Event acknowledged successfully",
	})
}

// GetEventSnapshot handles GET /api/v1/events/{id}/snapshot
func (h *EventHandler) GetEventSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		utils.RespondBadRequest(w, "Event ID is required", nil)
		return
	}

	// Get the event
	event, err := h.eventService.GetEvent(ctx, id)
	if err != nil {
		utils.RespondNotFound(w, "Event not found")
		return
	}

	// Check if event has a snapshot path
	if event.SnapshotPath == "" {
		// No snapshot stored, try to get a live snapshot from the camera
		if h.cameraService == nil {
			utils.RespondError(w, http.StatusNotFound, "NO_SNAPSHOT", "No snapshot available for this event", nil)
			return
		}

		// Get camera client
		client, err := h.cameraService.GetCameraClient(event.CameraID)
		if err != nil {
			utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found or unavailable", nil)
			return
		}

		// Get live snapshot from camera
		snapshot, err := client.GetSnapshot(ctx, 0) // Channel 0
		if err != nil {
			utils.RespondError(w, http.StatusInternalServerError, "SNAPSHOT_ERROR", "Failed to capture snapshot from camera", nil)
			return
		}

		// Return the snapshot image
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(snapshot)))
		w.WriteHeader(http.StatusOK)
		w.Write(snapshot)
		return
	}

	// Serve the stored snapshot file
	// Check if file exists
	if _, err := os.Stat(event.SnapshotPath); os.IsNotExist(err) {
		utils.RespondError(w, http.StatusNotFound, "SNAPSHOT_NOT_FOUND", "Snapshot file not found", nil)
		return
	}

	// Read the file
	data, err := os.ReadFile(event.SnapshotPath)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "FILE_READ_ERROR", "Failed to read snapshot file", nil)
		return
	}

	// Determine content type based on file extension
	contentType := "image/jpeg"
	ext := filepath.Ext(event.SnapshotPath)
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".bmp":
		contentType = "image/bmp"
	}

	// Return the snapshot image
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
