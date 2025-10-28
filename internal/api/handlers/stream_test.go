package handlers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	reolink "github.com/mosleyit/reolink_api_wrapper"
	"github.com/mosleyit/reolink_server/internal/api/service"
)

// MockStreamService is a mock implementation of StreamServiceInterface
type MockStreamService struct {
	mock.Mock
}

func (m *MockStreamService) ProxyFLVStream(ctx context.Context, cameraID string, streamType reolink.StreamType, channel int, w io.Writer) error {
	args := m.Called(ctx, cameraID, streamType, channel, w)
	return args.Error(0)
}

func (m *MockStreamService) StartHLSStream(ctx context.Context, cameraID string, streamType reolink.StreamType, channel int) (*service.StreamSession, error) {
	args := m.Called(ctx, cameraID, streamType, channel)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*service.StreamSession), args.Error(1)
}

func (m *MockStreamService) GetHLSPlaylist(sessionID string) (string, error) {
	args := m.Called(sessionID)
	return args.String(0), args.Error(1)
}

func (m *MockStreamService) GetHLSSegment(sessionID, segmentName string) (string, error) {
	args := m.Called(sessionID, segmentName)
	return args.String(0), args.Error(1)
}

func (m *MockStreamService) StopSession(sessionID string) error {
	args := m.Called(sessionID)
	return args.Error(0)
}

func TestNewStreamHandler(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.streamService)
}

func TestStreamHandler_ProxyFLV_MissingCameraID(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cameras//stream/flv/proxy", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.ProxyFLV(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Camera ID is required")
}

func TestStreamHandler_ProxyFLV_Success(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	mockService.On("ProxyFLVStream", mock.Anything, "cam-123", reolink.StreamMain, 0, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cameras/cam-123/stream/flv/proxy", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "cam-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.ProxyFLV(w, req)

	assert.Equal(t, "video/x-flv", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
	mockService.AssertExpectations(t)
}

func TestStreamHandler_ProxyFLV_WithStreamType(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	mockService.On("ProxyFLVStream", mock.Anything, "cam-123", reolink.StreamSub, 0, mock.Anything).Return(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cameras/cam-123/stream/flv/proxy?stream=sub", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "cam-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.ProxyFLV(w, req)

	mockService.AssertExpectations(t)
}

func TestStreamHandler_StartHLS_MissingCameraID(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cameras//stream/hls/start", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.StartHLS(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Camera ID is required")
}

func TestStreamHandler_StartHLS_Success(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	session := &service.StreamSession{
		ID:         "session-123",
		CameraID:   "cam-123",
		StreamType: service.StreamTypeHLS,
		StartedAt:  time.Now(),
		ExpiresAt:  time.Now().Add(30 * time.Minute),
	}

	mockService.On("StartHLSStream", mock.Anything, "cam-123", reolink.StreamMain, 0).Return(session, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cameras/cam-123/stream/hls/start", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "cam-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.StartHLS(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "session-123")
	assert.Contains(t, w.Body.String(), "playlist.m3u8")
	mockService.AssertExpectations(t)
}

func TestStreamHandler_StartHLS_ServiceError(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	mockService.On("StartHLSStream", mock.Anything, "cam-123", reolink.StreamMain, 0).Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/cameras/cam-123/stream/hls/start", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "cam-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.StartHLS(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Failed to start HLS stream")
	mockService.AssertExpectations(t)
}

func TestStreamHandler_GetHLSPlaylist_MissingSessionID(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stream/hls//playlist.m3u8", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetHLSPlaylist(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Session ID is required")
}

func TestStreamHandler_GetHLSPlaylist_SessionNotFound(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	mockService.On("GetHLSPlaylist", "session-999").Return("", assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stream/hls/session-999/playlist.m3u8", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("session_id", "session-999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetHLSPlaylist(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Session not found")
	mockService.AssertExpectations(t)
}

func TestStreamHandler_GetHLSPlaylist_FileNotReady(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	mockService.On("GetHLSPlaylist", "session-123").Return("/nonexistent/playlist.m3u8", nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stream/hls/session-123/playlist.m3u8", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("session_id", "session-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetHLSPlaylist(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Playlist not ready")
	mockService.AssertExpectations(t)
}

func TestStreamHandler_GetHLSPlaylist_Success(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	// Create a temporary playlist file
	tmpDir := t.TempDir()
	playlistPath := filepath.Join(tmpDir, "playlist.m3u8")
	playlistContent := "#EXTM3U\n#EXT-X-VERSION:3\n"
	err := os.WriteFile(playlistPath, []byte(playlistContent), 0644)
	assert.NoError(t, err)

	mockService.On("GetHLSPlaylist", "session-123").Return(playlistPath, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stream/hls/session-123/playlist.m3u8", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("session_id", "session-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetHLSPlaylist(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/vnd.apple.mpegurl", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Body.String(), "#EXTM3U")
	mockService.AssertExpectations(t)
}

func TestStreamHandler_GetHLSSegment_MissingParameters(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stream/hls//segment_001.ts", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("segment", "segment_001.ts")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetHLSSegment(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestStreamHandler_GetHLSSegment_InvalidSegmentName(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/stream/hls/session-123/invalid.txt", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("session_id", "session-123")
	rctx.URLParams.Add("segment", "invalid.txt")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetHLSSegment(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid segment name")
}

func TestStreamHandler_StopHLS_MissingSessionID(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/stream/hls/", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.StopHLS(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Session ID is required")
}

func TestStreamHandler_StopHLS_Success(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	mockService.On("StopSession", "session-123").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/stream/hls/session-123", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("session_id", "session-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.StopHLS(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Session stopped successfully")
	mockService.AssertExpectations(t)
}

func TestStreamHandler_StopHLS_SessionNotFound(t *testing.T) {
	mockService := new(MockStreamService)
	handler := NewStreamHandler(mockService)

	mockService.On("StopSession", "session-999").Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/stream/hls/session-999", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("session_id", "session-999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.StopHLS(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Session not found")
	mockService.AssertExpectations(t)
}

