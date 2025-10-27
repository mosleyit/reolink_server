package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	reolink "github.com/mosleyit/reolink_api_wrapper"
	"github.com/mosleyit/reolink_server/internal/api/service"
	"github.com/mosleyit/reolink_server/internal/logger"
	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/mosleyit/reolink_server/pkg/utils"
)

// CameraHandler handles camera-related HTTP requests
type CameraHandler struct {
	cameraService *service.CameraService
}

// NewCameraHandler creates a new camera handler
func NewCameraHandler(cameraService *service.CameraService) *CameraHandler {
	return &CameraHandler{
		cameraService: cameraService,
	}
}

// ListCameras handles GET /api/v1/cameras
func (h *CameraHandler) ListCameras(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cameras, err := h.cameraService.ListCameras(ctx)
	if err != nil {
		logger.Error("Failed to list cameras", zap.Error(err))
		utils.RespondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve cameras", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"cameras": cameras,
		"total":   len(cameras),
	})
}

// AddCamera handles POST /api/v1/cameras
func (h *CameraHandler) AddCamera(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req models.CreateCameraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}

	// Validate required fields
	if req.Name == "" || req.Host == "" || req.Username == "" || req.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "VALIDATION_ERROR", "Missing required fields", map[string]interface{}{
			"required": []string{"name", "host", "username", "password"},
		})
		return
	}

	// Set default port if not provided
	if req.Port == 0 {
		if req.UseHTTPS {
			req.Port = 443
		} else {
			req.Port = 80
		}
	}

	// Create camera model
	camera := &models.Camera{
		Name:       req.Name,
		Host:       req.Host,
		Port:       req.Port,
		Username:   req.Username,
		Password:   req.Password,
		UseHTTPS:   req.UseHTTPS,
		SkipVerify: req.SkipVerify,
		Status:     "offline",
	}

	// Add camera via service
	if err := h.cameraService.AddCamera(ctx, camera); err != nil {
		logger.Error("Failed to add camera", zap.Error(err), zap.String("name", req.Name))
		utils.RespondError(w, http.StatusInternalServerError, "ADD_CAMERA_ERROR", "Failed to add camera", nil)
		return
	}

	logger.Info("Camera added successfully", zap.String("id", camera.ID), zap.String("name", camera.Name))
	utils.RespondJSON(w, http.StatusCreated, camera)
}

// GetCamera handles GET /api/v1/cameras/{id}
func (h *CameraHandler) GetCamera(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	camera, err := h.cameraService.GetCamera(ctx, cameraID)
	if err != nil {
		logger.Error("Failed to get camera", zap.Error(err), zap.String("id", cameraID))
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, camera)
}

// UpdateCamera handles PUT /api/v1/cameras/{id}
func (h *CameraHandler) UpdateCamera(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	// Get existing camera
	camera, err := h.cameraService.GetCamera(ctx, cameraID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	var req models.UpdateCameraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}

	// Update fields if provided
	if req.Name != nil {
		camera.Name = *req.Name
	}
	if req.Host != nil {
		camera.Host = *req.Host
	}
	if req.Port != nil {
		camera.Port = *req.Port
	}
	if req.Username != nil {
		camera.Username = *req.Username
	}
	if req.Password != nil {
		camera.Password = *req.Password
	}
	if req.UseHTTPS != nil {
		camera.UseHTTPS = *req.UseHTTPS
	}
	if req.SkipVerify != nil {
		camera.SkipVerify = *req.SkipVerify
	}

	// Update camera via service
	if err := h.cameraService.UpdateCamera(ctx, camera); err != nil {
		logger.Error("Failed to update camera", zap.Error(err), zap.String("id", cameraID))
		utils.RespondError(w, http.StatusInternalServerError, "UPDATE_CAMERA_ERROR", "Failed to update camera", nil)
		return
	}

	logger.Info("Camera updated successfully", zap.String("id", cameraID))
	utils.RespondJSON(w, http.StatusOK, camera)
}

// DeleteCamera handles DELETE /api/v1/cameras/{id}
func (h *CameraHandler) DeleteCamera(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	if err := h.cameraService.DeleteCamera(ctx, cameraID); err != nil {
		logger.Error("Failed to delete camera", zap.Error(err), zap.String("id", cameraID))
		utils.RespondError(w, http.StatusInternalServerError, "DELETE_CAMERA_ERROR", "Failed to delete camera", nil)
		return
	}

	logger.Info("Camera deleted successfully", zap.String("id", cameraID))
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Camera deleted successfully",
		"id":      cameraID,
	})
}

// GetCameraStatus handles GET /api/v1/cameras/{id}/status
func (h *CameraHandler) GetCameraStatus(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	status, err := h.cameraService.GetCameraStatus(ctx, cameraID)
	if err != nil {
		logger.Error("Failed to get camera status", zap.Error(err), zap.String("id", cameraID))
		utils.RespondError(w, http.StatusInternalServerError, "STATUS_ERROR", "Failed to get camera status", nil)
		return
	}

	utils.RespondJSON(w, http.StatusOK, status)
}

// RebootCamera handles POST /api/v1/cameras/{id}/reboot
func (h *CameraHandler) RebootCamera(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	client, err := h.cameraService.GetCameraClient(cameraID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	if err := client.Reboot(ctx); err != nil {
		logger.Error("Failed to reboot camera", zap.Error(err), zap.String("id", cameraID))
		utils.RespondError(w, http.StatusInternalServerError, "REBOOT_ERROR", "Failed to reboot camera", nil)
		return
	}

	logger.Info("Camera rebooted successfully", zap.String("id", cameraID))
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Camera reboot initiated",
		"id":      cameraID,
	})
}

// GetSnapshot handles GET /api/v1/cameras/{id}/snapshot
func (h *CameraHandler) GetSnapshot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	client, err := h.cameraService.GetCameraClient(cameraID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	// Get snapshot from camera
	snapshot, err := client.GetSnapshot(ctx, 0) // Channel 0
	if err != nil {
		logger.Error("Failed to get snapshot", zap.Error(err), zap.String("id", cameraID))
		utils.RespondError(w, http.StatusInternalServerError, "SNAPSHOT_ERROR", "Failed to capture snapshot", nil)
		return
	}

	// Return image directly
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(snapshot)))
	w.WriteHeader(http.StatusOK)
	w.Write(snapshot)
}

// PTZMove handles POST /api/v1/cameras/{id}/ptz/move
func (h *CameraHandler) PTZMove(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	var req struct {
		Operation string `json:"operation"` // Up, Down, Left, Right, ZoomInc, ZoomDec, etc.
		Speed     *int   `json:"speed,omitempty"`
		Channel   *int   `json:"channel,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}

	client, err := h.cameraService.GetCameraClient(cameraID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	// Set defaults
	speed := 32
	if req.Speed != nil {
		speed = *req.Speed
	}
	channel := 0
	if req.Channel != nil {
		channel = *req.Channel
	}

	// Execute PTZ operation
	if err := client.PTZMove(ctx, req.Operation, speed, channel); err != nil {
		logger.Error("Failed to execute PTZ operation", zap.Error(err), zap.String("id", cameraID), zap.String("operation", req.Operation))
		utils.RespondError(w, http.StatusInternalServerError, "PTZ_ERROR", "Failed to execute PTZ operation", nil)
		return
	}

	logger.Info("PTZ operation executed", zap.String("id", cameraID), zap.String("operation", req.Operation))
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message":   "PTZ operation executed",
		"operation": req.Operation,
	})
}

// PTZPreset handles POST /api/v1/cameras/{id}/ptz/preset
func (h *CameraHandler) PTZPreset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	var req struct {
		PresetID int  `json:"preset_id"`
		Channel  *int `json:"channel,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}

	client, err := h.cameraService.GetCameraClient(cameraID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	channel := 0
	if req.Channel != nil {
		channel = *req.Channel
	}

	// Go to preset
	if err := client.PTZGotoPreset(ctx, channel, req.PresetID); err != nil {
		logger.Error("Failed to go to preset", zap.Error(err), zap.String("id", cameraID), zap.Int("preset", req.PresetID))
		utils.RespondError(w, http.StatusInternalServerError, "PTZ_ERROR", "Failed to go to preset", nil)
		return
	}

	logger.Info("Moved to PTZ preset", zap.String("id", cameraID), zap.Int("preset", req.PresetID))
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message":   "Moved to preset",
		"preset_id": req.PresetID,
	})
}

// ControlLED handles POST /api/v1/cameras/{id}/led
func (h *CameraHandler) ControlLED(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	var req struct {
		Type    string `json:"type"`    // "white" or "ir"
		Enabled bool   `json:"enabled"` // true or false
		Channel *int   `json:"channel,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}

	client, err := h.cameraService.GetCameraClient(cameraID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	channel := 0
	if req.Channel != nil {
		channel = *req.Channel
	}

	// Control LED based on type
	var controlErr error
	switch req.Type {
	case "white":
		// SetWhiteLED requires a WhiteLed config
		mode := 0
		if req.Enabled {
			mode = 1
		}
		config := &reolink.WhiteLed{
			Channel: channel,
			Mode:    mode,
		}
		controlErr = client.SetWhiteLED(ctx, config)
	case "ir":
		// SetIRLights requires a state string
		state := "Auto"
		if req.Enabled {
			state = "On"
		} else {
			state = "Off"
		}
		controlErr = client.SetIRLights(ctx, channel, state)
	default:
		utils.RespondError(w, http.StatusBadRequest, "INVALID_LED_TYPE", "LED type must be 'white' or 'ir'", nil)
		return
	}

	if controlErr != nil {
		logger.Error("Failed to control LED", zap.Error(controlErr), zap.String("id", cameraID), zap.String("type", req.Type))
		utils.RespondError(w, http.StatusInternalServerError, "LED_ERROR", "Failed to control LED", nil)
		return
	}

	logger.Info("LED controlled", zap.String("id", cameraID), zap.String("type", req.Type), zap.Bool("enabled", req.Enabled))
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "LED controlled successfully",
		"type":    req.Type,
		"enabled": req.Enabled,
	})
}

// TriggerSiren handles POST /api/v1/cameras/{id}/siren
func (h *CameraHandler) TriggerSiren(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	var req struct {
		Duration int  `json:"duration"` // Duration in seconds
		Channel  *int `json:"channel,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "INVALID_JSON", "Invalid request body", nil)
		return
	}

	client, err := h.cameraService.GetCameraClient(cameraID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	channel := 0
	if req.Channel != nil {
		channel = *req.Channel
	}

	// Trigger audio alarm (siren)
	if err := client.TriggerSiren(ctx, channel, req.Duration); err != nil {
		logger.Error("Failed to trigger siren", zap.Error(err), zap.String("id", cameraID))
		utils.RespondError(w, http.StatusInternalServerError, "SIREN_ERROR", "Failed to trigger siren", nil)
		return
	}

	logger.Info("Siren triggered", zap.String("id", cameraID), zap.Int("duration", req.Duration))
	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message":  "Siren triggered",
		"duration": req.Duration,
	})
}

// GetCameraConfig handles GET /api/v1/cameras/{id}/config/{type}
func (h *CameraHandler) GetCameraConfig(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement get camera config - this would require mapping config types to SDK methods
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Get camera config not yet implemented", nil)
}

// UpdateCameraConfig handles PUT /api/v1/cameras/{id}/config/{type}
func (h *CameraHandler) UpdateCameraConfig(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement update camera config - this would require mapping config types to SDK methods
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "Update camera config not yet implemented", nil)
}

// GetCameraEvents handles GET /api/v1/cameras/{id}/events
func (h *CameraHandler) GetCameraEvents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50 // default
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	// Get events from service
	events, err := h.cameraService.GetCameraEvents(ctx, cameraID, limit, offset)
	if err != nil {
		logger.Error("Failed to get camera events", zap.Error(err), zap.String("id", cameraID))
		utils.RespondError(w, http.StatusInternalServerError, "DATABASE_ERROR", "Failed to retrieve events", nil)
		return
	}

	// Get total count
	total, err := h.cameraService.CountCameraEvents(ctx, cameraID)
	if err != nil {
		logger.Error("Failed to count camera events", zap.Error(err), zap.String("id", cameraID))
		total = len(events) // fallback to current count
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"events": events,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// GetRTSPURL handles GET /api/v1/cameras/{id}/stream/rtsp
func (h *CameraHandler) GetRTSPURL(w http.ResponseWriter, r *http.Request) {
	cameraID := chi.URLParam(r, "id")

	client, err := h.cameraService.GetCameraClient(cameraID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	// Get stream type from query parameter (default to main stream)
	streamTypeStr := r.URL.Query().Get("stream")
	streamType := reolink.StreamMain
	switch streamTypeStr {
	case "sub":
		streamType = reolink.StreamSub
	case "ext":
		streamType = reolink.StreamExt
	}

	// Get channel from query parameter (default to 0)
	channelStr := r.URL.Query().Get("channel")
	channel := 0
	if channelStr != "" {
		if c, err := strconv.Atoi(channelStr); err == nil {
			channel = c
		}
	}

	rtspURL := client.GetRTSPURL(streamType, channel)

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"url":         rtspURL,
		"stream_type": streamTypeStr,
		"channel":     channel,
	})
}

// GetFLVURL handles GET /api/v1/cameras/{id}/stream/flv
func (h *CameraHandler) GetFLVURL(w http.ResponseWriter, r *http.Request) {
	cameraID := chi.URLParam(r, "id")

	client, err := h.cameraService.GetCameraClient(cameraID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "CAMERA_NOT_FOUND", "Camera not found", nil)
		return
	}

	// Get stream type from query parameter (default to main stream)
	streamTypeStr := r.URL.Query().Get("stream")
	streamType := reolink.StreamMain
	switch streamTypeStr {
	case "sub":
		streamType = reolink.StreamSub
	case "ext":
		streamType = reolink.StreamExt
	}

	// Get channel from query parameter (default to 0)
	channelStr := r.URL.Query().Get("channel")
	channel := 0
	if channelStr != "" {
		if c, err := strconv.Atoi(channelStr); err == nil {
			channel = c
		}
	}

	flvURL := client.GetFLVURL(streamType, channel)

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"url":         flvURL,
		"stream_type": streamTypeStr,
		"channel":     channel,
	})
}

// GetHLSURL handles GET /api/v1/cameras/{id}/stream/hls
func (h *CameraHandler) GetHLSURL(w http.ResponseWriter, r *http.Request) {
	// HLS would require transcoding - not implemented yet
	utils.RespondError(w, http.StatusNotImplemented, "NOT_IMPLEMENTED", "HLS streaming requires transcoding which is not yet implemented", nil)
}
