package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mosleyit/reolink_server/internal/api/service"
	"github.com/mosleyit/reolink_server/pkg/utils"
)

// EventHandler handles event-related HTTP requests
type EventHandler struct {
	eventService *service.EventService
}

// NewEventHandler creates a new event handler
func NewEventHandler(eventService *service.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
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
	// TODO: Implement get event snapshot
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Get event snapshot not yet implemented", nil)
}
