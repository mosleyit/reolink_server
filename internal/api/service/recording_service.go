package service

import (
	"context"
	"fmt"
	"time"

	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/storage/models"
)

// RecordingRepository interface for dependency injection
type RecordingRepository interface {
	GetByID(ctx context.Context, id string) (*models.Recording, error)
	ListByCameraID(ctx context.Context, cameraID string, limit, offset int) ([]*models.Recording, error)
	ListByTimeRange(ctx context.Context, cameraID string, startTime, endTime time.Time, limit, offset int) ([]*models.Recording, error)
	Search(ctx context.Context, req *models.RecordingSearchRequest) ([]*models.Recording, error)
	Count(ctx context.Context) (int, error)
	CountByCameraID(ctx context.Context, cameraID string) (int, error)
	GetTotalSize(ctx context.Context) (int64, error)
	GetTotalSizeByCameraID(ctx context.Context, cameraID string) (int64, error)
	Create(ctx context.Context, recording *models.Recording) error
	Delete(ctx context.Context, id string) error
	DeleteOlderThan(ctx context.Context, olderThan time.Time) (int64, error)
}

// CameraManager interface for dependency injection
type CameraManager interface {
	GetCamera(id string) (*camera.CameraClient, error)
}

// RecordingService handles recording operations
type RecordingService struct {
	recordingRepo RecordingRepository
	cameraManager CameraManager
}

// NewRecordingService creates a new recording service
func NewRecordingService(recordingRepo RecordingRepository, cameraManager CameraManager) *RecordingService {
	return &RecordingService{
		recordingRepo: recordingRepo,
		cameraManager: cameraManager,
	}
}

// GetRecording retrieves a recording by ID
func (s *RecordingService) GetRecording(ctx context.Context, id string) (*models.Recording, error) {
	return s.recordingRepo.GetByID(ctx, id)
}

// ListRecordings retrieves all recordings with pagination
func (s *RecordingService) ListRecordings(ctx context.Context, limit, offset int) ([]*models.Recording, error) {
	// Default limit
	if limit <= 0 {
		limit = 50
	}

	// Use Search with empty criteria to get all recordings
	req := &models.RecordingSearchRequest{
		Limit:  limit,
		Offset: offset,
	}

	return s.recordingRepo.Search(ctx, req)
}

// ListRecordingsByCameraID retrieves recordings for a specific camera
func (s *RecordingService) ListRecordingsByCameraID(ctx context.Context, cameraID string, limit, offset int) ([]*models.Recording, error) {
	if limit <= 0 {
		limit = 50
	}

	return s.recordingRepo.ListByCameraID(ctx, cameraID, limit, offset)
}

// ListRecordingsByTimeRange retrieves recordings within a time range for a specific camera
func (s *RecordingService) ListRecordingsByTimeRange(ctx context.Context, cameraID string, startTime, endTime time.Time, limit, offset int) ([]*models.Recording, error) {
	if limit <= 0 {
		limit = 50
	}

	return s.recordingRepo.ListByTimeRange(ctx, cameraID, startTime, endTime, limit, offset)
}

// SearchRecordings searches recordings with flexible criteria
func (s *RecordingService) SearchRecordings(ctx context.Context, req *models.RecordingSearchRequest) ([]*models.Recording, error) {
	// Set default limit
	if req.Limit <= 0 {
		req.Limit = 50
	}

	return s.recordingRepo.Search(ctx, req)
}

// CountRecordings returns the total count of recordings
func (s *RecordingService) CountRecordings(ctx context.Context) (int, error) {
	return s.recordingRepo.Count(ctx)
}

// CountRecordingsByCameraID returns the count of recordings for a specific camera
func (s *RecordingService) CountRecordingsByCameraID(ctx context.Context, cameraID string) (int, error) {
	return s.recordingRepo.CountByCameraID(ctx, cameraID)
}

// GetTotalSize returns the total size of all recordings in bytes
func (s *RecordingService) GetTotalSize(ctx context.Context) (int64, error) {
	return s.recordingRepo.GetTotalSize(ctx)
}

// GetTotalSizeByCameraID returns the total size of recordings for a specific camera
func (s *RecordingService) GetTotalSizeByCameraID(ctx context.Context, cameraID string) (int64, error) {
	return s.recordingRepo.GetTotalSizeByCameraID(ctx, cameraID)
}

// DeleteRecording deletes a recording by ID
func (s *RecordingService) DeleteRecording(ctx context.Context, id string) error {
	return s.recordingRepo.Delete(ctx, id)
}

// DeleteOldRecordings deletes recordings older than the specified time
func (s *RecordingService) DeleteOldRecordings(ctx context.Context, olderThan time.Time) (int64, error) {
	return s.recordingRepo.DeleteOlderThan(ctx, olderThan)
}

// CreateRecording creates a new recording
func (s *RecordingService) CreateRecording(ctx context.Context, recording *models.Recording) error {
	return s.recordingRepo.Create(ctx, recording)
}

// RecordingDownloadInfo contains information for downloading a recording
type RecordingDownloadInfo struct {
	Recording   *models.Recording `json:"recording"`
	DownloadURL string            `json:"download_url"`
	Method      string            `json:"method"`
	Note        string            `json:"note,omitempty"`
}

// GetRecordingDownloadInfo generates download information for a recording
func (s *RecordingService) GetRecordingDownloadInfo(ctx context.Context, recording *models.Recording) (*RecordingDownloadInfo, error) {
	// Get camera client to generate download URL
	cameraClient, err := s.cameraManager.GetCamera(recording.CameraID)
	if err != nil {
		return nil, fmt.Errorf("camera not found or unavailable: %w", err)
	}

	// Generate download URL using the SDK's Download method
	// The SDK's Download method returns a URL for downloading the recording
	downloadURL := cameraClient.Download(recording.StoragePath, "")

	return &RecordingDownloadInfo{
		Recording:   recording,
		DownloadURL: downloadURL,
		Method:      "GET",
		Note:        "Use this URL to download the recording file directly from the camera",
	}, nil
}
