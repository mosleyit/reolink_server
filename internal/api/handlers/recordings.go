package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/mosleyit/reolink_server/pkg/utils"
)

// RecordingServiceInterface defines the interface for recording service
type RecordingServiceInterface interface {
	GetRecording(ctx context.Context, id string) (*models.Recording, error)
	ListRecordings(ctx context.Context, limit, offset int) ([]*models.Recording, error)
	ListRecordingsByCameraID(ctx context.Context, cameraID string, limit, offset int) ([]*models.Recording, error)
	ListRecordingsByTimeRange(ctx context.Context, cameraID string, startTime, endTime time.Time, limit, offset int) ([]*models.Recording, error)
	SearchRecordings(ctx context.Context, req *models.RecordingSearchRequest) ([]*models.Recording, error)
	CountRecordings(ctx context.Context) (int, error)
	GetTotalSize(ctx context.Context) (int64, error)
	DeleteRecording(ctx context.Context, id string) error
}

// RecordingHandler handles recording requests
type RecordingHandler struct {
	recordingService RecordingServiceInterface
}

// NewRecordingHandler creates a new recording handler
func NewRecordingHandler(recordingService RecordingServiceInterface) *RecordingHandler {
	return &RecordingHandler{recordingService: recordingService}
}

// ListRecordings handles GET /api/v1/recordings
func (h *RecordingHandler) ListRecordings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	cameraID := r.URL.Query().Get("camera_id")
	startTimeStr := r.URL.Query().Get("start_time")
	endTimeStr := r.URL.Query().Get("end_time")

	var recordings []*models.Recording
	var err error

	// Filter by camera ID and/or time range
	if cameraID != "" && startTimeStr != "" && endTimeStr != "" {
		// Filter by both camera ID and time range
		startTime, err1 := time.Parse(time.RFC3339, startTimeStr)
		endTime, err2 := time.Parse(time.RFC3339, endTimeStr)

		if err1 != nil || err2 != nil {
			utils.RespondBadRequest(w, "Invalid time format. Use RFC3339 format (e.g., 2025-10-27T10:00:00Z)", nil)
			return
		}

		recordings, err = h.recordingService.ListRecordingsByTimeRange(ctx, cameraID, startTime, endTime, limit, offset)
	} else if cameraID != "" {
		// Filter by camera ID only
		recordings, err = h.recordingService.ListRecordingsByCameraID(ctx, cameraID, limit, offset)
	} else {
		// List all recordings
		recordings, err = h.recordingService.ListRecordings(ctx, limit, offset)
	}

	if err != nil {
		utils.RespondInternalError(w, "Failed to list recordings")
		return
	}

	// Get total count
	total, err := h.recordingService.CountRecordings(ctx)
	if err != nil {
		utils.RespondInternalError(w, "Failed to count recordings")
		return
	}

	// Get total size
	totalSize, err := h.recordingService.GetTotalSize(ctx)
	if err != nil {
		utils.RespondInternalError(w, "Failed to get total size")
		return
	}

	response := map[string]interface{}{
		"recordings": recordings,
		"total":      total,
		"total_size": totalSize,
		"limit":      limit,
		"offset":     offset,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// GetRecording handles GET /api/v1/recordings/{id}
func (h *RecordingHandler) GetRecording(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		utils.RespondBadRequest(w, "Recording ID is required", nil)
		return
	}

	recording, err := h.recordingService.GetRecording(ctx, id)
	if err != nil {
		utils.RespondNotFound(w, "Recording not found")
		return
	}

	utils.RespondJSON(w, http.StatusOK, recording)
}

// DownloadRecording handles GET /api/v1/recordings/{id}/download
func (h *RecordingHandler) DownloadRecording(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		utils.RespondBadRequest(w, "Recording ID is required", nil)
		return
	}

	recording, err := h.recordingService.GetRecording(ctx, id)
	if err != nil {
		utils.RespondNotFound(w, "Recording not found")
		return
	}

	// TODO: Implement actual file download from storage
	// For now, return recording metadata with download URL
	response := map[string]interface{}{
		"recording":    recording,
		"download_url": "/api/v1/recordings/" + id + "/file",
		"message":      "File download not yet implemented. Use the camera's recording download API directly.",
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// SearchRecordings handles POST /api/v1/recordings/search
func (h *RecordingHandler) SearchRecordings(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.RecordingSearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondBadRequest(w, "Invalid request body", nil)
		return
	}

	recordings, err := h.recordingService.SearchRecordings(ctx, &req)
	if err != nil {
		utils.RespondInternalError(w, "Failed to search recordings")
		return
	}

	// Get total count for the search
	total, err := h.recordingService.CountRecordings(ctx)
	if err != nil {
		utils.RespondInternalError(w, "Failed to count recordings")
		return
	}

	response := map[string]interface{}{
		"recordings": recordings,
		"total":      total,
		"limit":      req.Limit,
		"offset":     req.Offset,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// DeleteRecording handles DELETE /api/v1/recordings/{id}
func (h *RecordingHandler) DeleteRecording(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := chi.URLParam(r, "id")

	if id == "" {
		utils.RespondBadRequest(w, "Recording ID is required", nil)
		return
	}

	if err := h.recordingService.DeleteRecording(ctx, id); err != nil {
		utils.RespondInternalError(w, "Failed to delete recording")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "Recording deleted successfully",
	})
}
