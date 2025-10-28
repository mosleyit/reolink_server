package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/storage/models"
)

// MockRecordingRepository is a mock implementation of RecordingRepository
type MockRecordingRepository struct {
	mock.Mock
}

func (m *MockRecordingRepository) GetByID(ctx context.Context, id string) (*models.Recording, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Recording), args.Error(1)
}

func (m *MockRecordingRepository) ListByCameraID(ctx context.Context, cameraID string, limit, offset int) ([]*models.Recording, error) {
	args := m.Called(ctx, cameraID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Recording), args.Error(1)
}

func (m *MockRecordingRepository) ListByTimeRange(ctx context.Context, cameraID string, startTime, endTime time.Time, limit, offset int) ([]*models.Recording, error) {
	args := m.Called(ctx, cameraID, startTime, endTime, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Recording), args.Error(1)
}

func (m *MockRecordingRepository) Search(ctx context.Context, req *models.RecordingSearchRequest) ([]*models.Recording, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Recording), args.Error(1)
}

func (m *MockRecordingRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockRecordingRepository) CountByCameraID(ctx context.Context, cameraID string) (int, error) {
	args := m.Called(ctx, cameraID)
	return args.Int(0), args.Error(1)
}

func (m *MockRecordingRepository) GetTotalSize(ctx context.Context) (int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRecordingRepository) GetTotalSizeByCameraID(ctx context.Context, cameraID string) (int64, error) {
	args := m.Called(ctx, cameraID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockRecordingRepository) Create(ctx context.Context, recording *models.Recording) error {
	args := m.Called(ctx, recording)
	return args.Error(0)
}

func (m *MockRecordingRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRecordingRepository) DeleteOlderThan(ctx context.Context, olderThan time.Time) (int64, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).(int64), args.Error(1)
}

// MockCameraManager is a mock implementation of CameraManager
type MockCameraManager struct {
	mock.Mock
}

func (m *MockCameraManager) GetCamera(id string) (*camera.CameraClient, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*camera.CameraClient), args.Error(1)
}

func TestNewRecordingService(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.recordingRepo)
	assert.Equal(t, mockCameraManager, service.cameraManager)
}

func TestRecordingService_GetRecording(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	expectedRecording := &models.Recording{
		ID:       "rec-123",
		CameraID: "cam-123",
		FileName: "recording.mp4",
	}

	mockRepo.On("GetByID", ctx, "rec-123").Return(expectedRecording, nil)

	recording, err := service.GetRecording(ctx, "rec-123")

	assert.NoError(t, err)
	assert.Equal(t, expectedRecording, recording)
	mockRepo.AssertExpectations(t)
}

func TestRecordingService_ListRecordings(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	expectedRecordings := []*models.Recording{
		{ID: "rec-1", CameraID: "cam-1"},
		{ID: "rec-2", CameraID: "cam-2"},
	}

	searchReq := &models.RecordingSearchRequest{
		Limit:  50,
		Offset: 0,
	}

	mockRepo.On("Search", ctx, searchReq).Return(expectedRecordings, nil)

	recordings, err := service.ListRecordings(ctx, 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecordings, recordings)
	mockRepo.AssertExpectations(t)
}

func TestRecordingService_ListRecordingsByCameraID(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	expectedRecordings := []*models.Recording{
		{ID: "rec-1", CameraID: "cam-123"},
		{ID: "rec-2", CameraID: "cam-123"},
	}

	mockRepo.On("ListByCameraID", ctx, "cam-123", 50, 0).Return(expectedRecordings, nil)

	recordings, err := service.ListRecordingsByCameraID(ctx, "cam-123", 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecordings, recordings)
	mockRepo.AssertExpectations(t)
}

func TestRecordingService_ListRecordingsByTimeRange(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	startTime := time.Now().Add(-24 * time.Hour)
	endTime := time.Now()

	expectedRecordings := []*models.Recording{
		{ID: "rec-1", CameraID: "cam-123"},
	}

	mockRepo.On("ListByTimeRange", ctx, "cam-123", startTime, endTime, 50, 0).Return(expectedRecordings, nil)

	recordings, err := service.ListRecordingsByTimeRange(ctx, "cam-123", startTime, endTime, 0, 0)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecordings, recordings)
	mockRepo.AssertExpectations(t)
}

func TestRecordingService_SearchRecordings(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	cameraID := "cam-123"
	searchReq := &models.RecordingSearchRequest{
		CameraID: &cameraID,
		Limit:    10,
		Offset:   0,
	}

	expectedRecordings := []*models.Recording{
		{ID: "rec-1", CameraID: "cam-123"},
	}

	mockRepo.On("Search", ctx, searchReq).Return(expectedRecordings, nil)

	recordings, err := service.SearchRecordings(ctx, searchReq)

	assert.NoError(t, err)
	assert.Equal(t, expectedRecordings, recordings)
	mockRepo.AssertExpectations(t)
}

func TestRecordingService_CountRecordings(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	mockRepo.On("Count", ctx).Return(42, nil)

	count, err := service.CountRecordings(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 42, count)
	mockRepo.AssertExpectations(t)
}

func TestRecordingService_GetTotalSize(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	mockRepo.On("GetTotalSize", ctx).Return(int64(1024000), nil)

	size, err := service.GetTotalSize(ctx)

	assert.NoError(t, err)
	assert.Equal(t, int64(1024000), size)
	mockRepo.AssertExpectations(t)
}

func TestRecordingService_DeleteRecording(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	mockRepo.On("Delete", ctx, "rec-123").Return(nil)

	err := service.DeleteRecording(ctx, "rec-123")

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRecordingService_DeleteOldRecordings(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	olderThan := time.Now().Add(-30 * 24 * time.Hour)
	mockRepo.On("DeleteOlderThan", ctx, olderThan).Return(int64(5), nil)

	count, err := service.DeleteOldRecordings(ctx, olderThan)

	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	mockRepo.AssertExpectations(t)
}

func TestRecordingService_CreateRecording(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	recording := &models.Recording{
		ID:       "rec-123",
		CameraID: "cam-123",
		FileName: "recording.mp4",
	}

	mockRepo.On("Create", ctx, recording).Return(nil)

	err := service.CreateRecording(ctx, recording)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// Note: Testing GetRecordingDownloadInfo with a real camera client requires
// a fully initialized CameraClient with SDK client, which is complex to mock.
// The error case (camera not found) is tested below, and the success case
// is better tested at the integration level or through handler tests.

func TestRecordingService_GetRecordingDownloadInfo_CameraNotFound(t *testing.T) {
	mockRepo := new(MockRecordingRepository)
	mockCameraManager := new(MockCameraManager)
	service := NewRecordingService(mockRepo, mockCameraManager)
	ctx := context.Background()

	recording := &models.Recording{
		ID:          "rec-123",
		CameraID:    "cam-999",
		FileName:    "recording.mp4",
		StoragePath: "/path/to/recording.mp4",
	}

	mockCameraManager.On("GetCamera", "cam-999").Return(nil, assert.AnError)

	downloadInfo, err := service.GetRecordingDownloadInfo(ctx, recording)

	assert.Error(t, err)
	assert.Nil(t, downloadInfo)
	assert.Contains(t, err.Error(), "camera not found or unavailable")
	mockCameraManager.AssertExpectations(t)
}
