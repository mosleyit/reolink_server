package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/storage/models"
)

// MockEventService is a mock implementation of EventService
type MockEventService struct {
	mock.Mock
}

func (m *MockEventService) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventService) ListEvents(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockEventService) CountEvents(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockEventService) AcknowledgeEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockCameraServiceForEvents is a mock implementation of CameraServiceInterface for event tests
type MockCameraServiceForEvents struct {
	mock.Mock
}

func (m *MockCameraServiceForEvents) AddCamera(ctx context.Context, camera *models.Camera) error {
	args := m.Called(ctx, camera)
	return args.Error(0)
}

func (m *MockCameraServiceForEvents) GetCamera(ctx context.Context, id string) (*models.Camera, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Camera), args.Error(1)
}

func (m *MockCameraServiceForEvents) ListCameras(ctx context.Context) ([]*models.Camera, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Camera), args.Error(1)
}

func (m *MockCameraServiceForEvents) UpdateCamera(ctx context.Context, camera *models.Camera) error {
	args := m.Called(ctx, camera)
	return args.Error(0)
}

func (m *MockCameraServiceForEvents) DeleteCamera(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCameraServiceForEvents) GetCameraStatus(ctx context.Context, id string) (*models.CameraStatus, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CameraStatus), args.Error(1)
}

func (m *MockCameraServiceForEvents) GetCameraClient(id string) (*camera.CameraClient, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*camera.CameraClient), args.Error(1)
}

func (m *MockCameraServiceForEvents) GetCameraEvents(ctx context.Context, cameraID string, limit, offset int) ([]*models.Event, error) {
	args := m.Called(ctx, cameraID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockCameraServiceForEvents) CountCameraEvents(ctx context.Context, cameraID string) (int, error) {
	args := m.Called(ctx, cameraID)
	return args.Int(0), args.Error(1)
}

func TestNewEventHandler(t *testing.T) {
	mockEventService := new(MockEventService)
	mockCameraService := new(MockCameraServiceForEvents)
	handler := NewEventHandler(mockEventService, mockCameraService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockEventService, handler.eventService)
	assert.Equal(t, mockCameraService, handler.cameraService)
}

func TestEventHandler_ListEvents(t *testing.T) {
	mockEventService := new(MockEventService)
	mockCameraService := new(MockCameraServiceForEvents)
	handler := NewEventHandler(mockEventService, mockCameraService)

	events := []*models.Event{
		{ID: "evt-1", CameraID: "cam-1", Type: models.EventMotionDetected},
		{ID: "evt-2", CameraID: "cam-2", Type: models.EventAIPerson},
	}

	mockEventService.On("ListEvents", mock.Anything, 50, 0).Return(events, nil)
	mockEventService.On("CountEvents", mock.Anything).Return(2, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events", nil)
	w := httptest.NewRecorder()

	handler.ListEvents(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "evt-1")
	assert.Contains(t, w.Body.String(), "evt-2")
	mockEventService.AssertExpectations(t)
}

func TestEventHandler_GetEvent(t *testing.T) {
	mockEventService := new(MockEventService)
	mockCameraService := new(MockCameraServiceForEvents)
	handler := NewEventHandler(mockEventService, mockCameraService)

	event := &models.Event{
		ID:       "evt-123",
		CameraID: "cam-123",
		Type:     models.EventMotionDetected,
	}

	mockEventService.On("GetEvent", mock.Anything, "evt-123").Return(event, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/evt-123", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "evt-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetEvent(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "evt-123")
	mockEventService.AssertExpectations(t)
}

func TestEventHandler_AcknowledgeEvent(t *testing.T) {
	mockEventService := new(MockEventService)
	mockCameraService := new(MockCameraServiceForEvents)
	handler := NewEventHandler(mockEventService, mockCameraService)

	mockEventService.On("AcknowledgeEvent", mock.Anything, "evt-123").Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/events/evt-123/acknowledge", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "evt-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.AcknowledgeEvent(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "acknowledged successfully")
	mockEventService.AssertExpectations(t)
}

func TestEventHandler_GetEventSnapshot_WithStoredSnapshot(t *testing.T) {
	mockEventService := new(MockEventService)
	mockCameraService := new(MockCameraServiceForEvents)
	handler := NewEventHandler(mockEventService, mockCameraService)

	// Create a temporary snapshot file
	tmpDir := t.TempDir()
	snapshotPath := filepath.Join(tmpDir, "snapshot.jpg")
	snapshotData := []byte("fake jpeg data")
	err := os.WriteFile(snapshotPath, snapshotData, 0644)
	assert.NoError(t, err)

	event := &models.Event{
		ID:           "evt-123",
		CameraID:     "cam-123",
		Type:         models.EventMotionDetected,
		SnapshotPath: snapshotPath,
	}

	mockEventService.On("GetEvent", mock.Anything, "evt-123").Return(event, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/evt-123/snapshot", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "evt-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetEventSnapshot(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "image/jpeg", w.Header().Get("Content-Type"))
	assert.Equal(t, snapshotData, w.Body.Bytes())
	mockEventService.AssertExpectations(t)
}

func TestEventHandler_GetEventSnapshot_MissingID(t *testing.T) {
	mockEventService := new(MockEventService)
	mockCameraService := new(MockCameraServiceForEvents)
	handler := NewEventHandler(mockEventService, mockCameraService)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events//snapshot", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetEventSnapshot(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Event ID is required")
}

func TestEventHandler_GetEventSnapshot_EventNotFound(t *testing.T) {
	mockEventService := new(MockEventService)
	mockCameraService := new(MockCameraServiceForEvents)
	handler := NewEventHandler(mockEventService, mockCameraService)

	mockEventService.On("GetEvent", mock.Anything, "evt-999").Return(nil, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/evt-999/snapshot", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "evt-999")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetEventSnapshot(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Event not found")
	mockEventService.AssertExpectations(t)
}

func TestEventHandler_GetEventSnapshot_NoSnapshotNoCameraService(t *testing.T) {
	mockEventService := new(MockEventService)
	handler := NewEventHandler(mockEventService, nil) // No camera service

	event := &models.Event{
		ID:           "evt-123",
		CameraID:     "cam-123",
		Type:         models.EventMotionDetected,
		SnapshotPath: "", // No snapshot
	}

	mockEventService.On("GetEvent", mock.Anything, "evt-123").Return(event, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/evt-123/snapshot", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "evt-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetEventSnapshot(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "No snapshot available")
	mockEventService.AssertExpectations(t)
}

func TestEventHandler_GetEventSnapshot_SnapshotFileNotFound(t *testing.T) {
	mockEventService := new(MockEventService)
	mockCameraService := new(MockCameraServiceForEvents)
	handler := NewEventHandler(mockEventService, mockCameraService)

	event := &models.Event{
		ID:           "evt-123",
		CameraID:     "cam-123",
		Type:         models.EventMotionDetected,
		SnapshotPath: "/nonexistent/path/snapshot.jpg",
	}

	mockEventService.On("GetEvent", mock.Anything, "evt-123").Return(event, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/events/evt-123/snapshot", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "evt-123")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetEventSnapshot(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Snapshot file not found")
	mockEventService.AssertExpectations(t)
}
