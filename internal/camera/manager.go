package camera

import (
	"context"
	"fmt"
	"sync"
	"time"

	reolink "github.com/mosleyit/reolink_api_wrapper"
	"go.uber.org/zap"

	"github.com/mosleyit/reolink_server/internal/logger"
	"github.com/mosleyit/reolink_server/internal/storage/models"
)

// CameraRepository interface for database operations
type CameraRepository interface {
	UpdateStatus(ctx context.Context, id string, status string, lastSeen time.Time) error
}

// Manager manages all camera connections and operations
type Manager struct {
	cameras map[string]*CameraClient
	mu      sync.RWMutex
	config  *Config
	repo    CameraRepository
}

// Config holds camera manager configuration
type Config struct {
	HealthCheckInterval time.Duration
	ConnectionTimeout   time.Duration
	MaxRetries          int
	RetryBackoff        time.Duration
}

// CameraClient wraps a Reolink API client with additional metadata
type CameraClient struct {
	Camera       *models.Camera
	Client       *reolink.Client
	LastHealthy  time.Time
	FailureCount int
	CircuitOpen  bool
	mu           sync.RWMutex
}

// NewManager creates a new camera manager
func NewManager(config *Config, repo CameraRepository) *Manager {
	if config == nil {
		config = &Config{
			HealthCheckInterval: 30 * time.Second,
			ConnectionTimeout:   10 * time.Second,
			MaxRetries:          3,
			RetryBackoff:        5 * time.Second,
		}
	}

	return &Manager{
		cameras: make(map[string]*CameraClient),
		config:  config,
		repo:    repo,
	}
}

// AddCamera adds a new camera to the manager
func (m *Manager) AddCamera(ctx context.Context, camera *models.Camera) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Validate input
	if camera == nil {
		return fmt.Errorf("camera cannot be nil")
	}
	if camera.ID == "" {
		return fmt.Errorf("camera ID cannot be empty")
	}
	if camera.Host == "" {
		return fmt.Errorf("camera host cannot be empty")
	}
	if camera.Port <= 0 {
		return fmt.Errorf("camera port must be positive")
	}
	if camera.Username == "" {
		return fmt.Errorf("camera username cannot be empty")
	}

	// Check if camera already exists
	if _, exists := m.cameras[camera.ID]; exists {
		return fmt.Errorf("camera %s already exists", camera.ID)
	}

	// Create Reolink client
	client, err := m.createClient(camera)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Test connection
	if err := client.Login(ctx); err != nil {
		return fmt.Errorf("failed to connect to camera: %w", err)
	}

	// Get device info
	info, err := client.System.GetDeviceInfo(ctx)
	if err != nil {
		logger.Warn("Failed to get device info", zap.String("camera_id", camera.ID), zap.Error(err))
	} else {
		camera.Model = info.Model
		camera.FirmwareVer = info.FirmVer
		camera.HardwareVer = info.HardVer // SDK uses HardVer field
	}

	// Add to manager
	m.cameras[camera.ID] = &CameraClient{
		Camera:      camera,
		Client:      client,
		LastHealthy: time.Now(),
	}

	camera.Status = "online"
	camera.LastSeen = time.Now()

	logger.Info("Camera added",
		zap.String("camera_id", camera.ID),
		zap.String("name", camera.Name),
		zap.String("model", camera.Model),
	)

	return nil
}

// RemoveCamera removes a camera from the manager
func (m *Manager) RemoveCamera(cameraID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	client, exists := m.cameras[cameraID]
	if !exists {
		return fmt.Errorf("camera %s not found", cameraID)
	}

	// Logout from camera
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = client.Client.Logout(ctx)

	delete(m.cameras, cameraID)

	logger.Info("Camera removed", zap.String("camera_id", cameraID))
	return nil
}

// GetCamera returns a camera client by ID
func (m *Manager) GetCamera(cameraID string) (*CameraClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.cameras[cameraID]
	if !exists {
		return nil, fmt.Errorf("camera %s not found", cameraID)
	}

	return client, nil
}

// ListCameras returns all cameras
func (m *Manager) ListCameras() []*models.Camera {
	m.mu.RLock()
	defer m.mu.RUnlock()

	cameras := make([]*models.Camera, 0, len(m.cameras))
	for _, client := range m.cameras {
		cameras = append(cameras, client.Camera)
	}

	return cameras
}

// GetCameraStatus returns the current status of a camera
func (m *Manager) GetCameraStatus(cameraID string) (*models.CameraStatus, error) {
	client, err := m.GetCamera(cameraID)
	if err != nil {
		return nil, err
	}

	client.mu.RLock()
	defer client.mu.RUnlock()

	status := &models.CameraStatus{
		CameraID:    client.Camera.ID,
		Status:      client.Camera.Status,
		Model:       client.Camera.Model,
		FirmwareVer: client.Camera.FirmwareVer,
		LastSeen:    client.Camera.LastSeen,
	}

	return status, nil
}

// HealthCheck performs health checks on all cameras
func (m *Manager) HealthCheck(ctx context.Context) {
	m.mu.RLock()
	cameras := make([]*CameraClient, 0, len(m.cameras))
	for _, client := range m.cameras {
		cameras = append(cameras, client)
	}
	m.mu.RUnlock()

	for _, client := range cameras {
		go m.checkCameraHealth(ctx, client)
	}
}

// checkCameraHealth checks the health of a single camera
func (m *Manager) checkCameraHealth(ctx context.Context, client *CameraClient) {
	client.mu.Lock()
	defer client.mu.Unlock()

	// Skip if circuit is open
	if client.CircuitOpen {
		logger.Debug("Circuit open, skipping health check",
			zap.String("camera_id", client.Camera.ID),
		)
		return
	}

	// Perform health check
	_, err := client.Client.System.GetDeviceInfo(ctx)
	if err != nil {
		client.FailureCount++
		oldStatus := client.Camera.Status
		client.Camera.Status = "offline"

		logger.Warn("Camera health check failed",
			zap.String("camera_id", client.Camera.ID),
			zap.Int("failure_count", client.FailureCount),
			zap.Error(err),
		)

		// Update database if status changed
		if m.repo != nil && oldStatus != "offline" {
			if err := m.repo.UpdateStatus(ctx, client.Camera.ID, "offline", time.Now()); err != nil {
				logger.Error("Failed to update camera status in database",
					zap.String("camera_id", client.Camera.ID),
					zap.Error(err))
			}
		}

		// Open circuit if too many failures
		if client.FailureCount >= m.config.MaxRetries {
			client.CircuitOpen = true
			logger.Error("Circuit opened for camera",
				zap.String("camera_id", client.Camera.ID),
			)
		}
	} else {
		// Reset on success
		client.FailureCount = 0
		client.CircuitOpen = false
		client.LastHealthy = time.Now()
		oldStatus := client.Camera.Status
		client.Camera.Status = "online"
		client.Camera.LastSeen = time.Now()

		// Update database if status changed
		if m.repo != nil && oldStatus != "online" {
			if err := m.repo.UpdateStatus(ctx, client.Camera.ID, "online", time.Now()); err != nil {
				logger.Error("Failed to update camera status in database",
					zap.String("camera_id", client.Camera.ID),
					zap.Error(err))
			}
		}
	}
}

// StartHealthMonitoring starts periodic health monitoring
func (m *Manager) StartHealthMonitoring(ctx context.Context) {
	ticker := time.NewTicker(m.config.HealthCheckInterval)
	defer ticker.Stop()

	logger.Info("Health monitoring started",
		zap.Duration("interval", m.config.HealthCheckInterval),
	)

	for {
		select {
		case <-ctx.Done():
			logger.Info("Health monitoring stopped")
			return
		case <-ticker.C:
			m.HealthCheck(ctx)
		}
	}
}

// createClient creates a new Reolink API client
func (m *Manager) createClient(camera *models.Camera) (*reolink.Client, error) {
	opts := []reolink.Option{
		reolink.WithCredentials(camera.Username, camera.Password),
		reolink.WithTimeout(m.config.ConnectionTimeout),
	}

	if camera.UseHTTPS {
		opts = append(opts, reolink.WithHTTPS(true))
	}

	if camera.SkipVerify {
		opts = append(opts, reolink.WithInsecureSkipVerify(true))
	}

	host := camera.Host
	if camera.Port > 0 && camera.Port != 80 && camera.Port != 443 {
		host = fmt.Sprintf("%s:%d", camera.Host, camera.Port)
	}

	client := reolink.NewClient(host, opts...)
	return client, nil
}

// Shutdown gracefully shuts down the manager
func (m *Manager) Shutdown(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	logger.Info("Shutting down camera manager")

	for id, client := range m.cameras {
		if err := client.Client.Logout(ctx); err != nil {
			logger.Warn("Failed to logout from camera",
				zap.String("camera_id", id),
				zap.Error(err),
			)
		}
	}

	m.cameras = make(map[string]*CameraClient)
	return nil
}
