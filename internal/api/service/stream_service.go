package service

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	reolink "github.com/mosleyit/reolink_api_wrapper"
	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/logger"
)

// StreamType represents the type of stream
type StreamType string

const (
	StreamTypeFLV  StreamType = "flv"
	StreamTypeHLS  StreamType = "hls"
	StreamTypeRTSP StreamType = "rtsp"
	StreamTypeRTMP StreamType = "rtmp"
)

// StreamSession represents an active streaming session
type StreamSession struct {
	ID         string
	CameraID   string
	StreamType StreamType
	StartedAt  time.Time
	LastAccess time.Time
	ExpiresAt  time.Time
	cancel     context.CancelFunc
}

// CameraManagerInterface defines the interface for camera manager operations
type CameraManagerInterface interface {
	GetCamera(id string) (*camera.CameraClient, error)
}

// StreamService manages video streaming sessions
type StreamService struct {
	cameraManager CameraManagerInterface
	sessions      map[string]*StreamSession
	sessionsMu    sync.RWMutex
	hlsOutputDir  string
	ffmpegPath    string
}

// StreamServiceConfig holds configuration for the stream service
type StreamServiceConfig struct {
	HLSOutputDir    string
	FFmpegPath      string
	SessionTimeout  time.Duration
	CleanupInterval time.Duration
}

// NewStreamService creates a new stream service
func NewStreamService(cameraManager CameraManagerInterface, config *StreamServiceConfig) *StreamService {
	if config == nil {
		config = &StreamServiceConfig{
			HLSOutputDir:    "/tmp/hls",
			FFmpegPath:      "ffmpeg",
			SessionTimeout:  30 * time.Minute,
			CleanupInterval: 5 * time.Minute,
		}
	}

	// Create HLS output directory if it doesn't exist
	if err := os.MkdirAll(config.HLSOutputDir, 0755); err != nil {
		logger.Error("Failed to create HLS output directory", zap.Error(err))
	}

	service := &StreamService{
		cameraManager: cameraManager,
		sessions:      make(map[string]*StreamSession),
		hlsOutputDir:  config.HLSOutputDir,
		ffmpegPath:    config.FFmpegPath,
	}

	// Start session cleanup goroutine
	go service.cleanupExpiredSessions(config.CleanupInterval)

	return service
}

// ProxyFLVStream proxies an FLV stream from the camera
func (s *StreamService) ProxyFLVStream(ctx context.Context, cameraID string, streamType reolink.StreamType, channel int, w io.Writer) error {
	client, err := s.cameraManager.GetCamera(cameraID)
	if err != nil {
		return fmt.Errorf("camera not found: %w", err)
	}

	// Get FLV URL from camera
	flvURL := client.GetFLVURL(streamType, channel)
	if flvURL == "" {
		return fmt.Errorf("failed to get FLV URL for camera %s", cameraID)
	}

	logger.Info("Proxying FLV stream",
		zap.String("camera_id", cameraID),
		zap.String("url", flvURL))

	// Create HTTP request to camera's FLV stream
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, flvURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to camera stream: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("camera returned status %d", resp.StatusCode)
	}

	// Copy stream data to response writer
	_, err = io.Copy(w, resp.Body)
	if err != nil && err != io.EOF && ctx.Err() == nil {
		return fmt.Errorf("failed to proxy stream: %w", err)
	}

	return nil
}

// StartHLSStream starts an HLS transcoding session
func (s *StreamService) StartHLSStream(ctx context.Context, cameraID string, streamType reolink.StreamType, channel int) (*StreamSession, error) {
	client, err := s.cameraManager.GetCamera(cameraID)
	if err != nil {
		return nil, fmt.Errorf("camera not found: %w", err)
	}

	// Get RTSP URL from camera (HLS transcoding uses RTSP as input)
	rtspURL := client.GetRTSPURL(streamType, channel)
	if rtspURL == "" {
		return nil, fmt.Errorf("failed to get RTSP URL for camera %s", cameraID)
	}

	// Create session
	sessionID := uuid.New().String()
	sessionDir := filepath.Join(s.hlsOutputDir, sessionID)
	if err := os.MkdirAll(sessionDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory: %w", err)
	}

	playlistPath := filepath.Join(sessionDir, "playlist.m3u8")

	// Create context for FFmpeg process
	ffmpegCtx, cancel := context.WithCancel(ctx)

	// Build FFmpeg command
	// ffmpeg -i rtsp://camera/stream -c:v copy -c:a aac -f hls \
	//        -hls_time 2 -hls_list_size 5 -hls_flags delete_segments \
	//        -hls_segment_filename 'segment_%03d.ts' playlist.m3u8
	cmd := exec.CommandContext(ffmpegCtx, s.ffmpegPath,
		"-i", rtspURL,
		"-c:v", "copy",
		"-c:a", "aac",
		"-f", "hls",
		"-hls_time", "2",
		"-hls_list_size", "5",
		"-hls_flags", "delete_segments",
		"-hls_segment_filename", filepath.Join(sessionDir, "segment_%03d.ts"),
		playlistPath,
	)

	// Capture stderr for error logging
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start FFmpeg process
	if err := cmd.Start(); err != nil {
		cancel()
		return nil, fmt.Errorf("failed to start FFmpeg: %w", err)
	}

	logger.Info("Started HLS transcoding session",
		zap.String("session_id", sessionID),
		zap.String("camera_id", cameraID),
		zap.String("rtsp_url", rtspURL))

	// Create session
	session := &StreamSession{
		ID:         sessionID,
		CameraID:   cameraID,
		StreamType: StreamTypeHLS,
		StartedAt:  time.Now(),
		LastAccess: time.Now(),
		ExpiresAt:  time.Now().Add(30 * time.Minute),
		cancel:     cancel,
	}

	// Store session
	s.sessionsMu.Lock()
	s.sessions[sessionID] = session
	s.sessionsMu.Unlock()

	// Read FFmpeg stderr in background
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := stderrPipe.Read(buf)
			if n > 0 {
				logger.Debug("FFmpeg output",
					zap.String("session_id", sessionID),
					zap.String("output", string(buf[:n])))
			}
			if err != nil {
				break
			}
		}
	}()

	// Monitor FFmpeg process
	go func() {
		err := cmd.Wait()
		if err != nil && ffmpegCtx.Err() == nil {
			logger.Error("FFmpeg process exited with error",
				zap.String("session_id", sessionID),
				zap.Error(err))
		}

		// Cleanup session
		s.StopSession(sessionID)
	}()

	return session, nil
}

// GetHLSPlaylist returns the path to the HLS playlist for a session
func (s *StreamService) GetHLSPlaylist(sessionID string) (string, error) {
	s.sessionsMu.RLock()
	session, exists := s.sessions[sessionID]
	s.sessionsMu.RUnlock()

	if !exists {
		return "", fmt.Errorf("session not found")
	}

	// Update last access time
	s.sessionsMu.Lock()
	session.LastAccess = time.Now()
	session.ExpiresAt = time.Now().Add(30 * time.Minute)
	s.sessionsMu.Unlock()

	playlistPath := filepath.Join(s.hlsOutputDir, sessionID, "playlist.m3u8")
	return playlistPath, nil
}

// GetHLSSegment returns the path to an HLS segment for a session
func (s *StreamService) GetHLSSegment(sessionID, segmentName string) (string, error) {
	s.sessionsMu.RLock()
	session, exists := s.sessions[sessionID]
	s.sessionsMu.RUnlock()

	if !exists {
		return "", fmt.Errorf("session not found")
	}

	// Update last access time
	s.sessionsMu.Lock()
	session.LastAccess = time.Now()
	session.ExpiresAt = time.Now().Add(30 * time.Minute)
	s.sessionsMu.Unlock()

	segmentPath := filepath.Join(s.hlsOutputDir, sessionID, segmentName)

	// Validate segment path to prevent directory traversal
	if !filepath.HasPrefix(segmentPath, filepath.Join(s.hlsOutputDir, sessionID)) {
		return "", fmt.Errorf("invalid segment path")
	}

	return segmentPath, nil
}

// StopSession stops a streaming session
func (s *StreamService) StopSession(sessionID string) error {
	s.sessionsMu.Lock()
	session, exists := s.sessions[sessionID]
	if !exists {
		s.sessionsMu.Unlock()
		return fmt.Errorf("session not found")
	}
	delete(s.sessions, sessionID)
	s.sessionsMu.Unlock()

	// Cancel context to stop FFmpeg
	if session.cancel != nil {
		session.cancel()
	}

	// Cleanup session directory
	sessionDir := filepath.Join(s.hlsOutputDir, sessionID)
	if err := os.RemoveAll(sessionDir); err != nil {
		logger.Error("Failed to cleanup session directory",
			zap.String("session_id", sessionID),
			zap.Error(err))
	}

	logger.Info("Stopped streaming session",
		zap.String("session_id", sessionID),
		zap.String("camera_id", session.CameraID))

	return nil
}

// cleanupExpiredSessions periodically cleans up expired sessions
func (s *StreamService) cleanupExpiredSessions(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		s.sessionsMu.Lock()
		for sessionID, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				logger.Info("Cleaning up expired session",
					zap.String("session_id", sessionID),
					zap.String("camera_id", session.CameraID))

				// Cancel context
				if session.cancel != nil {
					session.cancel()
				}

				// Remove from map
				delete(s.sessions, sessionID)

				// Cleanup directory
				sessionDir := filepath.Join(s.hlsOutputDir, sessionID)
				os.RemoveAll(sessionDir)
			}
		}
		s.sessionsMu.Unlock()
	}
}
