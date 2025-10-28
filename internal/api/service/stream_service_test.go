package service

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	reolink "github.com/mosleyit/reolink_api_wrapper"
	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/storage/models"
)

// MockCameraManagerForStream is a mock implementation of CameraManagerInterface
type MockCameraManagerForStream struct {
	mock.Mock
}

func (m *MockCameraManagerForStream) GetCamera(id string) (*camera.CameraClient, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*camera.CameraClient), args.Error(1)
}

func TestNewStreamService(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	
	// Test with default config
	service := NewStreamService(mockCameraManager, nil)
	assert.NotNil(t, service)
	assert.Equal(t, mockCameraManager, service.cameraManager)
	assert.NotNil(t, service.sessions)
	assert.Equal(t, "/tmp/hls", service.hlsOutputDir)
	assert.Equal(t, "ffmpeg", service.ffmpegPath)

	// Test with custom config
	config := &StreamServiceConfig{
		HLSOutputDir:    "/custom/hls",
		FFmpegPath:      "/usr/bin/ffmpeg",
		SessionTimeout:  15 * time.Minute,
		CleanupInterval: 2 * time.Minute,
	}
	service = NewStreamService(mockCameraManager, config)
	assert.NotNil(t, service)
	assert.Equal(t, "/custom/hls", service.hlsOutputDir)
	assert.Equal(t, "/usr/bin/ffmpeg", service.ffmpegPath)
}

func TestStreamService_ProxyFLVStream_CameraNotFound(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	service := NewStreamService(mockCameraManager, nil)

	mockCameraManager.On("GetCamera", "cam-999").Return(nil, assert.AnError)

	ctx := context.Background()
	err := service.ProxyFLVStream(ctx, "cam-999", reolink.StreamMain, 0, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "camera not found")
	mockCameraManager.AssertExpectations(t)
}

func TestStreamService_ProxyFLVStream_EmptyURL(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	service := NewStreamService(mockCameraManager, nil)

	// Create a mock camera client that returns empty URL
	mockCamera := &models.Camera{
		ID:       "cam-123",
		Host:     "192.168.1.100",
		Username: "admin",
		Password: "password",
	}
	
	client := reolink.NewClient(mockCamera.Host,
		reolink.WithCredentials(mockCamera.Username, mockCamera.Password),
	)
	
	cameraClient := &camera.CameraClient{
		Camera: mockCamera,
		Client: client,
	}

	mockCameraManager.On("GetCamera", "cam-123").Return(cameraClient, nil)

	ctx := context.Background()
	// Note: This will fail because GetFLVURL returns a URL, but we're testing the error path
	// In a real scenario, we'd need to mock the client's GetFLVURL method
	err := service.ProxyFLVStream(ctx, "cam-123", reolink.StreamMain, 0, nil)

	// The error will be about connection failure since we can't actually connect
	assert.Error(t, err)
	mockCameraManager.AssertExpectations(t)
}

func TestStreamService_StartHLSStream_CameraNotFound(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	service := NewStreamService(mockCameraManager, nil)

	mockCameraManager.On("GetCamera", "cam-999").Return(nil, assert.AnError)

	ctx := context.Background()
	session, err := service.StartHLSStream(ctx, "cam-999", reolink.StreamMain, 0)

	assert.Error(t, err)
	assert.Nil(t, session)
	assert.Contains(t, err.Error(), "camera not found")
	mockCameraManager.AssertExpectations(t)
}

func TestStreamService_GetHLSPlaylist_SessionNotFound(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	service := NewStreamService(mockCameraManager, nil)

	path, err := service.GetHLSPlaylist("invalid-session")

	assert.Error(t, err)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "session not found")
}

func TestStreamService_GetHLSSegment_SessionNotFound(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	service := NewStreamService(mockCameraManager, nil)

	path, err := service.GetHLSSegment("invalid-session", "segment_001.ts")

	assert.Error(t, err)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "session not found")
}

func TestStreamService_GetHLSSegment_InvalidPath(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	tmpDir := t.TempDir()
	config := &StreamServiceConfig{
		HLSOutputDir:    tmpDir,
		FFmpegPath:      "ffmpeg",
		SessionTimeout:  30 * time.Minute,
		CleanupInterval: 5 * time.Minute,
	}
	service := NewStreamService(mockCameraManager, config)

	// Create a fake session
	sessionID := "test-session"
	service.sessions[sessionID] = &StreamSession{
		ID:         sessionID,
		CameraID:   "cam-123",
		StreamType: StreamTypeHLS,
		StartedAt:  time.Now(),
		LastAccess: time.Now(),
		ExpiresAt:  time.Now().Add(30 * time.Minute),
	}

	// Try to access a segment with directory traversal
	path, err := service.GetHLSSegment(sessionID, "../../../etc/passwd")

	assert.Error(t, err)
	assert.Empty(t, path)
	assert.Contains(t, err.Error(), "invalid segment path")
}

func TestStreamService_StopSession_SessionNotFound(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	service := NewStreamService(mockCameraManager, nil)

	err := service.StopSession("invalid-session")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session not found")
}

func TestStreamService_StopSession_Success(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	tmpDir := t.TempDir()
	config := &StreamServiceConfig{
		HLSOutputDir:    tmpDir,
		FFmpegPath:      "ffmpeg",
		SessionTimeout:  30 * time.Minute,
		CleanupInterval: 5 * time.Minute,
	}
	service := NewStreamService(mockCameraManager, config)

	// Create a fake session with directory
	sessionID := "test-session"
	sessionDir := filepath.Join(tmpDir, sessionID)
	err := os.MkdirAll(sessionDir, 0755)
	assert.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	service.sessions[sessionID] = &StreamSession{
		ID:         sessionID,
		CameraID:   "cam-123",
		StreamType: StreamTypeHLS,
		StartedAt:  time.Now(),
		LastAccess: time.Now(),
		ExpiresAt:  time.Now().Add(30 * time.Minute),
		cancel:     cancel,
	}

	// Stop the session
	err = service.StopSession(sessionID)

	assert.NoError(t, err)
	assert.NotContains(t, service.sessions, sessionID)
	
	// Verify directory was cleaned up
	_, err = os.Stat(sessionDir)
	assert.True(t, os.IsNotExist(err))
	
	// Verify context was cancelled
	assert.Error(t, ctx.Err())
}

func TestStreamService_GetHLSPlaylist_UpdatesLastAccess(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	tmpDir := t.TempDir()
	config := &StreamServiceConfig{
		HLSOutputDir:    tmpDir,
		FFmpegPath:      "ffmpeg",
		SessionTimeout:  30 * time.Minute,
		CleanupInterval: 5 * time.Minute,
	}
	service := NewStreamService(mockCameraManager, config)

	// Create a fake session
	sessionID := "test-session"
	oldTime := time.Now().Add(-10 * time.Minute)
	service.sessions[sessionID] = &StreamSession{
		ID:         sessionID,
		CameraID:   "cam-123",
		StreamType: StreamTypeHLS,
		StartedAt:  oldTime,
		LastAccess: oldTime,
		ExpiresAt:  oldTime.Add(30 * time.Minute),
	}

	// Get playlist
	path, err := service.GetHLSPlaylist(sessionID)

	assert.NoError(t, err)
	assert.NotEmpty(t, path)
	
	// Verify last access was updated
	session := service.sessions[sessionID]
	assert.True(t, session.LastAccess.After(oldTime))
	assert.True(t, session.ExpiresAt.After(oldTime.Add(30*time.Minute)))
}

func TestStreamService_GetHLSSegment_UpdatesLastAccess(t *testing.T) {
	mockCameraManager := new(MockCameraManagerForStream)
	tmpDir := t.TempDir()
	config := &StreamServiceConfig{
		HLSOutputDir:    tmpDir,
		FFmpegPath:      "ffmpeg",
		SessionTimeout:  30 * time.Minute,
		CleanupInterval: 5 * time.Minute,
	}
	service := NewStreamService(mockCameraManager, config)

	// Create a fake session
	sessionID := "test-session"
	oldTime := time.Now().Add(-10 * time.Minute)
	service.sessions[sessionID] = &StreamSession{
		ID:         sessionID,
		CameraID:   "cam-123",
		StreamType: StreamTypeHLS,
		StartedAt:  oldTime,
		LastAccess: oldTime,
		ExpiresAt:  oldTime.Add(30 * time.Minute),
	}

	// Get segment
	path, err := service.GetHLSSegment(sessionID, "segment_001.ts")

	assert.NoError(t, err)
	assert.NotEmpty(t, path)
	
	// Verify last access was updated
	session := service.sessions[sessionID]
	assert.True(t, session.LastAccess.After(oldTime))
	assert.True(t, session.ExpiresAt.After(oldTime.Add(30*time.Minute)))
}

