package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockCameraServiceForConfig is a mock for CameraService focusing on config operations
type MockCameraServiceForConfig struct {
	mock.Mock
}

func (m *MockCameraServiceForConfig) AddCamera(ctx context.Context, camera *models.Camera) error {
	args := m.Called(ctx, camera)
	return args.Error(0)
}

func (m *MockCameraServiceForConfig) GetCamera(ctx context.Context, id string) (*models.Camera, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Camera), args.Error(1)
}

func (m *MockCameraServiceForConfig) ListCameras(ctx context.Context) ([]*models.Camera, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Camera), args.Error(1)
}

func (m *MockCameraServiceForConfig) UpdateCamera(ctx context.Context, camera *models.Camera) error {
	args := m.Called(ctx, camera)
	return args.Error(0)
}

func (m *MockCameraServiceForConfig) DeleteCamera(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCameraServiceForConfig) GetCameraStatus(ctx context.Context, id string) (*models.CameraStatus, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CameraStatus), args.Error(1)
}

func (m *MockCameraServiceForConfig) GetCameraClient(id string) (*camera.CameraClient, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*camera.CameraClient), args.Error(1)
}

func (m *MockCameraServiceForConfig) GetCameraEvents(ctx context.Context, cameraID string, limit, offset int) ([]*models.Event, error) {
	args := m.Called(ctx, cameraID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Event), args.Error(1)
}

func (m *MockCameraServiceForConfig) CountCameraEvents(ctx context.Context, cameraID string) (int, error) {
	args := m.Called(ctx, cameraID)
	return args.Int(0), args.Error(1)
}

func TestCameraHandler_GetCameraConfig_DeviceName(t *testing.T) {
	mockService := new(MockCameraServiceForConfig)
	handler := &CameraHandler{cameraService: mockService}

	// Create a mock camera client
	// Note: We can't easily mock the camera client methods without a full mock implementation
	// For now, we'll test the error paths

	mockService.On("GetCameraClient", "camera-123").Return(nil, errors.New("camera not found"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cameras/camera-123/config/device_name", nil)
	w := httptest.NewRecorder()

	// Add URL params using chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "camera-123")
	rctx.URLParams.Add("type", "device_name")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetCameraConfig(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestCameraHandler_GetCameraConfig_UnsupportedType(t *testing.T) {
	mockService := new(MockCameraServiceForConfig)
	handler := &CameraHandler{cameraService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cameras/camera-123/config/invalid_type", nil)
	w := httptest.NewRecorder()

	// Add URL params using chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "camera-123")
	rctx.URLParams.Add("type", "invalid_type")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetCameraConfig(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Error   map[string]interface{} `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error["details"].(map[string]interface{}), "supported_types")
}

func TestCameraHandler_GetCameraConfig_InvalidChannel(t *testing.T) {
	mockService := new(MockCameraServiceForConfig)
	handler := &CameraHandler{cameraService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cameras/camera-123/config/encoding?channel=invalid", nil)
	w := httptest.NewRecorder()

	// Add URL params using chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "camera-123")
	rctx.URLParams.Add("type", "encoding")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetCameraConfig(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Error   map[string]interface{} `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Equal(t, "BAD_REQUEST", response.Error["code"])
}

func TestCameraHandler_UpdateCameraConfig_DeviceName(t *testing.T) {
	mockService := new(MockCameraServiceForConfig)
	handler := &CameraHandler{cameraService: mockService}

	mockService.On("GetCameraClient", "camera-123").Return(nil, errors.New("camera not found"))

	reqBody := map[string]string{"name": "New Camera Name"}
	bodyBytes, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/cameras/camera-123/config/device_name", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add URL params using chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "camera-123")
	rctx.URLParams.Add("type", "device_name")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.UpdateCameraConfig(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestCameraHandler_UpdateCameraConfig_UnsupportedType(t *testing.T) {
	mockService := new(MockCameraServiceForConfig)
	handler := &CameraHandler{cameraService: mockService}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/cameras/camera-123/config/unsupported", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add URL params using chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "camera-123")
	rctx.URLParams.Add("type", "unsupported")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.UpdateCameraConfig(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Error   map[string]interface{} `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
	assert.Contains(t, response.Error["details"].(map[string]interface{}), "supported_types")
}

func TestCameraHandler_UpdateCameraConfig_InvalidJSON(t *testing.T) {
	mockService := new(MockCameraServiceForConfig)
	handler := &CameraHandler{cameraService: mockService}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/cameras/camera-123/config/device_name", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add URL params using chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "camera-123")
	rctx.URLParams.Add("type", "device_name")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.UpdateCameraConfig(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response struct {
		Success bool                   `json:"success"`
		Error   map[string]interface{} `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.False(t, response.Success)
}

func TestCameraHandler_UpdateCameraConfig_TimeConfig_InvalidJSON(t *testing.T) {
	mockService := new(MockCameraServiceForConfig)
	handler := &CameraHandler{cameraService: mockService}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/cameras/camera-123/config/time", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add URL params using chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "camera-123")
	rctx.URLParams.Add("type", "time")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.UpdateCameraConfig(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCameraHandler_UpdateCameraConfig_SystemConfig_InvalidJSON(t *testing.T) {
	mockService := new(MockCameraServiceForConfig)
	handler := &CameraHandler{cameraService: mockService}

	req := httptest.NewRequest(http.MethodPut, "/api/v1/cameras/camera-123/config/system", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Add URL params using chi context
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "camera-123")
	rctx.URLParams.Add("type", "system")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.UpdateCameraConfig(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Test that all supported config types are properly listed
func TestCameraHandler_GetCameraConfig_SupportedTypes(t *testing.T) {
	supportedTypes := []string{
		"time", "device_name", "auto_maint", "system", "encoding", "ai",
		"motion_alarm", "alarm", "audio_alarm", "buzzer_alarm", "ai_alarm",
		"recording", "osd", "image", "isp", "mask", "crop",
		"network_port", "ntp", "wifi", "email", "ftp", "push",
	}

	mockService := new(MockCameraServiceForConfig)
	handler := &CameraHandler{cameraService: mockService}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/cameras/camera-123/config/invalid", nil)
	w := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "camera-123")
	rctx.URLParams.Add("type", "invalid")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetCameraConfig(w, req)

	var response struct {
		Success bool                   `json:"success"`
		Error   map[string]interface{} `json:"error"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	details := response.Error["details"].(map[string]interface{})
	returnedTypes := details["supported_types"].([]interface{})

	assert.Equal(t, len(supportedTypes), len(returnedTypes))
}
