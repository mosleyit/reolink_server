package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mosleyit/reolink_server/internal/storage/models"
)

// MockRecordingService is a mock implementation of RecordingService
type MockRecordingService struct {
	mock.Mock
}

func (m *MockRecordingService) GetRecording(ctx context.Context, id string) (*models.Recording, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Recording), args.Error(1)
}

func (m *MockRecordingService) ListRecordings(ctx context.Context, limit, offset int) ([]*models.Recording, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Recording), args.Error(1)
}

func (m *MockRecordingService) ListRecordingsByCameraID(ctx context.Context, cameraID string, limit, offset int) ([]*models.Recording, error) {
	args := m.Called(ctx, cameraID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Recording), args.Error(1)
}

func (m *MockRecordingService) ListRecordingsByTimeRange(ctx context.Context, cameraID string, startTime, endTime time.Time, limit, offset int) ([]*models.Recording, error) {
	args := m.Called(ctx, cameraID, startTime, endTime, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Recording), args.Error(1)
}

func (m *MockRecordingService) SearchRecordings(ctx context.Context, req *models.RecordingSearchRequest) ([]*models.Recording, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Recording), args.Error(1)
}

func (m *MockRecordingService) CountRecordings(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockRecordingService) CountRecordingsByCameraID(ctx context.Context, cameraID string) (int, error) {
	args := m.Called(ctx, cameraID)
	return args.Int(0), args.Error(1)
}

func (m *MockRecordingService) GetTotalSize(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRecordingService) GetTotalSizeByCameraID(ctx context.Context, cameraID string) (int64, error) {
	args := m.Called(ctx, cameraID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRecordingService) DeleteRecording(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRecordingService) DeleteOldRecordings(ctx context.Context, olderThan time.Time) (int64, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRecordingService) CreateRecording(ctx context.Context, recording *models.Recording) error {
	args := m.Called(ctx, recording)
	return args.Error(0)
}

func TestNewRecordingHandler(t *testing.T) {
	mockService := new(MockRecordingService)
	handler := NewRecordingHandler(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.recordingService)
}

func TestRecordingHandler_ListRecordings(t *testing.T) {
	mockService := new(MockRecordingService)
	handler := NewRecordingHandler(mockService)

	recordings := []*models.Recording{
		{ID: "rec-1", CameraID: "cam-1", FileName: "recording1.mp4"},
		{ID: "rec-2", CameraID: "cam-2", FileName: "recording2.mp4"},
	}

	mockService.On("ListRecordings", mock.Anything, mock.Anything, mock.Anything).Return(recordings, nil)
	mockService.On("CountRecordings", mock.Anything).Return(2, nil)
	mockService.On("GetTotalSize", mock.Anything).Return(int64(2048000), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recordings", nil)
	w := httptest.NewRecorder()

	handler.ListRecordings(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "rec-1")
	assert.Contains(t, w.Body.String(), "rec-2")
	mockService.AssertExpectations(t)
}

func TestRecordingHandler_ListRecordingsByCameraID(t *testing.T) {
	mockService := new(MockRecordingService)
	handler := NewRecordingHandler(mockService)

	recordings := []*models.Recording{
		{ID: "rec-1", CameraID: "cam-123", FileName: "recording1.mp4"},
	}

	mockService.On("ListRecordingsByCameraID", mock.Anything, "cam-123", mock.Anything, mock.Anything).Return(recordings, nil)
	mockService.On("CountRecordings", mock.Anything).Return(1, nil)
	mockService.On("GetTotalSize", mock.Anything).Return(int64(1024000), nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recordings?camera_id=cam-123", nil)
	w := httptest.NewRecorder()

	handler.ListRecordings(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "cam-123")
	mockService.AssertExpectations(t)
}

func TestRecordingHandler_GetRecording(t *testing.T) {
	mockService := new(MockRecordingService)
	handler := NewRecordingHandler(mockService)

	recording := &models.Recording{
		ID:       "rec-123",
		CameraID: "cam-123",
		FileName: "recording.mp4",
	}

	mockService.On("GetRecording", mock.Anything, "rec-123").Return(recording, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recordings/rec-123", nil)
	w := httptest.NewRecorder()

	// Add URL params using chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "rec-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetRecording(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "rec-123")
	mockService.AssertExpectations(t)
}

func TestRecordingHandler_GetRecording_NotFound(t *testing.T) {
	mockService := new(MockRecordingService)
	handler := NewRecordingHandler(mockService)

	mockService.On("GetRecording", mock.Anything, "nonexistent").Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recordings/nonexistent", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "nonexistent")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetRecording(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestRecordingHandler_SearchRecordings(t *testing.T) {
	mockService := new(MockRecordingService)
	handler := NewRecordingHandler(mockService)

	cameraID := "cam-123"
	searchReq := &models.RecordingSearchRequest{
		CameraID: &cameraID,
		Limit:    10,
		Offset:   0,
	}

	recordings := []*models.Recording{
		{ID: "rec-1", CameraID: "cam-123"},
	}

	mockService.On("SearchRecordings", mock.Anything, mock.Anything).Return(recordings, nil)
	mockService.On("CountRecordings", mock.Anything).Return(1, nil)

	body, _ := json.Marshal(searchReq)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/recordings/search", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SearchRecordings(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "rec-1")
	mockService.AssertExpectations(t)
}

func TestRecordingHandler_SearchRecordings_InvalidJSON(t *testing.T) {
	mockService := new(MockRecordingService)
	handler := NewRecordingHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/recordings/search", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SearchRecordings(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecordingHandler_DeleteRecording(t *testing.T) {
	mockService := new(MockRecordingService)
	handler := NewRecordingHandler(mockService)

	mockService.On("DeleteRecording", mock.Anything, "rec-123").Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/recordings/rec-123", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "rec-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.DeleteRecording(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "deleted successfully")
	mockService.AssertExpectations(t)
}

func TestRecordingHandler_DownloadRecording(t *testing.T) {
	mockService := new(MockRecordingService)
	handler := NewRecordingHandler(mockService)

	recording := &models.Recording{
		ID:       "rec-123",
		CameraID: "cam-123",
		FileName: "recording.mp4",
	}

	mockService.On("GetRecording", mock.Anything, "rec-123").Return(recording, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/recordings/rec-123/download", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "rec-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.DownloadRecording(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "download_url")
	mockService.AssertExpectations(t)
}
