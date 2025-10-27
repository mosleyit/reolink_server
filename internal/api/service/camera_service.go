package service

import (
	"context"
	"fmt"

	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/mosleyit/reolink_server/internal/storage/repository"
)

// CameraService coordinates camera operations between the camera manager and database
type CameraService struct {
	cameraManager *camera.Manager
	cameraRepo    *repository.CameraRepository
	eventRepo     *repository.EventRepository
}

// NewCameraService creates a new camera service
func NewCameraService(
	cameraManager *camera.Manager,
	cameraRepo *repository.CameraRepository,
	eventRepo *repository.EventRepository,
) *CameraService {
	return &CameraService{
		cameraManager: cameraManager,
		cameraRepo:    cameraRepo,
		eventRepo:     eventRepo,
	}
}

// AddCamera adds a new camera to both the database and camera manager
func (s *CameraService) AddCamera(ctx context.Context, camera *models.Camera) error {
	// Save to database first
	if err := s.cameraRepo.Create(ctx, camera); err != nil {
		return fmt.Errorf("failed to save camera to database: %w", err)
	}

	// Add to camera manager
	if err := s.cameraManager.AddCamera(ctx, camera); err != nil {
		// Rollback: delete from database
		_ = s.cameraRepo.Delete(ctx, camera.ID)
		return fmt.Errorf("failed to add camera to manager: %w", err)
	}

	return nil
}

// GetCamera retrieves a camera by ID
func (s *CameraService) GetCamera(ctx context.Context, id string) (*models.Camera, error) {
	return s.cameraRepo.GetByID(ctx, id)
}

// ListCameras retrieves all cameras
func (s *CameraService) ListCameras(ctx context.Context) ([]*models.Camera, error) {
	return s.cameraRepo.List(ctx)
}

// UpdateCamera updates a camera in both database and manager
func (s *CameraService) UpdateCamera(ctx context.Context, camera *models.Camera) error {
	// Update in database
	if err := s.cameraRepo.Update(ctx, camera); err != nil {
		return fmt.Errorf("failed to update camera in database: %w", err)
	}

	// Remove from manager and re-add with new config
	s.cameraManager.RemoveCamera(camera.ID)
	if err := s.cameraManager.AddCamera(ctx, camera); err != nil {
		return fmt.Errorf("failed to update camera in manager: %w", err)
	}

	return nil
}

// DeleteCamera removes a camera from both database and manager
func (s *CameraService) DeleteCamera(ctx context.Context, id string) error {
	// Remove from manager first
	s.cameraManager.RemoveCamera(id)

	// Delete from database
	if err := s.cameraRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete camera from database: %w", err)
	}

	return nil
}

// GetCameraStatus retrieves the current status of a camera
func (s *CameraService) GetCameraStatus(ctx context.Context, id string) (*models.CameraStatus, error) {
	// Use the manager's GetCameraStatus method which already handles this
	return s.cameraManager.GetCameraStatus(id)
}

// GetCameraClient retrieves the camera client for direct SDK operations
func (s *CameraService) GetCameraClient(id string) (*camera.CameraClient, error) {
	return s.cameraManager.GetCamera(id)
}

// GetCameraEvents retrieves events for a specific camera
func (s *CameraService) GetCameraEvents(ctx context.Context, cameraID string, limit, offset int) ([]*models.Event, error) {
	return s.eventRepo.ListByCameraID(ctx, cameraID, limit, offset)
}

// CountCameraEvents returns the total number of events for a camera
func (s *CameraService) CountCameraEvents(ctx context.Context, cameraID string) (int, error) {
	return s.eventRepo.CountByCameraID(ctx, cameraID)
}
