package handlers

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	reolink "github.com/mosleyit/reolink_api_wrapper"
	"github.com/mosleyit/reolink_server/internal/api/service"
	"github.com/mosleyit/reolink_server/internal/logger"
	"github.com/mosleyit/reolink_server/pkg/utils"
)

// StreamServiceInterface defines the interface for stream service operations
type StreamServiceInterface interface {
	ProxyFLVStream(ctx context.Context, cameraID string, streamType reolink.StreamType, channel int, w io.Writer) error
	StartHLSStream(ctx context.Context, cameraID string, streamType reolink.StreamType, channel int) (*service.StreamSession, error)
	GetHLSPlaylist(sessionID string) (string, error)
	GetHLSSegment(sessionID, segmentName string) (string, error)
	StopSession(sessionID string) error
}

// StreamHandler handles streaming-related HTTP requests
type StreamHandler struct {
	streamService StreamServiceInterface
}

// NewStreamHandler creates a new stream handler
func NewStreamHandler(streamService StreamServiceInterface) *StreamHandler {
	return &StreamHandler{
		streamService: streamService,
	}
}

// ProxyFLV handles GET /api/v1/cameras/{id}/stream/flv/proxy
func (h *StreamHandler) ProxyFLV(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	if cameraID == "" {
		utils.RespondBadRequest(w, "Camera ID is required", nil)
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

	logger.Info("Starting FLV proxy stream",
		zap.String("camera_id", cameraID),
		zap.String("stream_type", streamTypeStr),
		zap.Int("channel", channel))

	// Set headers for FLV streaming
	w.Header().Set("Content-Type", "video/x-flv")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Connection", "keep-alive")

	// Proxy the stream
	if err := h.streamService.ProxyFLVStream(ctx, cameraID, streamType, channel, w); err != nil {
		logger.Error("Failed to proxy FLV stream",
			zap.String("camera_id", cameraID),
			zap.Error(err))
		// Can't send error response here as headers are already sent
		return
	}
}

// StartHLS handles POST /api/v1/cameras/{id}/stream/hls
func (h *StreamHandler) StartHLS(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cameraID := chi.URLParam(r, "id")

	if cameraID == "" {
		utils.RespondBadRequest(w, "Camera ID is required", nil)
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

	logger.Info("Starting HLS transcoding session",
		zap.String("camera_id", cameraID),
		zap.String("stream_type", streamTypeStr),
		zap.Int("channel", channel))

	// Start HLS session
	session, err := h.streamService.StartHLSStream(ctx, cameraID, streamType, channel)
	if err != nil {
		logger.Error("Failed to start HLS session",
			zap.String("camera_id", cameraID),
			zap.Error(err))
		utils.RespondInternalError(w, "Failed to start HLS stream")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"session_id":   session.ID,
		"camera_id":    session.CameraID,
		"stream_type":  streamTypeStr,
		"channel":      channel,
		"playlist_url": "/api/v1/stream/hls/" + session.ID + "/playlist.m3u8",
		"started_at":   session.StartedAt,
		"expires_at":   session.ExpiresAt,
	})
}

// GetHLSPlaylist handles GET /api/v1/stream/hls/{session_id}/playlist.m3u8
func (h *StreamHandler) GetHLSPlaylist(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "session_id")

	if sessionID == "" {
		utils.RespondBadRequest(w, "Session ID is required", nil)
		return
	}

	playlistPath, err := h.streamService.GetHLSPlaylist(sessionID)
	if err != nil {
		logger.Error("Failed to get HLS playlist",
			zap.String("session_id", sessionID),
			zap.Error(err))
		utils.RespondNotFound(w, "Session not found or expired")
		return
	}

	// Check if playlist file exists
	if _, err := os.Stat(playlistPath); os.IsNotExist(err) {
		utils.RespondNotFound(w, "Playlist not ready yet")
		return
	}

	// Serve the playlist file
	w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeFile(w, r, playlistPath)
}

// GetHLSSegment handles GET /api/v1/stream/hls/{session_id}/{segment}
func (h *StreamHandler) GetHLSSegment(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "session_id")
	segmentName := chi.URLParam(r, "segment")

	if sessionID == "" || segmentName == "" {
		utils.RespondBadRequest(w, "Session ID and segment name are required", nil)
		return
	}

	// Validate segment name (should be .ts file)
	if filepath.Ext(segmentName) != ".ts" {
		utils.RespondBadRequest(w, "Invalid segment name", nil)
		return
	}

	segmentPath, err := h.streamService.GetHLSSegment(sessionID, segmentName)
	if err != nil {
		logger.Error("Failed to get HLS segment",
			zap.String("session_id", sessionID),
			zap.String("segment", segmentName),
			zap.Error(err))
		utils.RespondNotFound(w, "Segment not found")
		return
	}

	// Check if segment file exists
	if _, err := os.Stat(segmentPath); os.IsNotExist(err) {
		utils.RespondNotFound(w, "Segment not found")
		return
	}

	// Serve the segment file
	w.Header().Set("Content-Type", "video/mp2t")
	w.Header().Set("Cache-Control", "public, max-age=31536000")
	http.ServeFile(w, r, segmentPath)
}

// StopHLS handles DELETE /api/v1/stream/hls/{session_id}
func (h *StreamHandler) StopHLS(w http.ResponseWriter, r *http.Request) {
	sessionID := chi.URLParam(r, "session_id")

	if sessionID == "" {
		utils.RespondBadRequest(w, "Session ID is required", nil)
		return
	}

	if err := h.streamService.StopSession(sessionID); err != nil {
		logger.Error("Failed to stop HLS session",
			zap.String("session_id", sessionID),
			zap.Error(err))
		utils.RespondNotFound(w, "Session not found")
		return
	}

	logger.Info("Stopped HLS session",
		zap.String("session_id", sessionID))

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Session stopped successfully",
	})
}
