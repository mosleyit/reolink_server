package camera

import (
	"context"
	"testing"
	"time"

	reolink "github.com/mosleyit/reolink_api_wrapper"
	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/stretchr/testify/assert"
)

func TestNewManager(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   *Config
	}{
		{
			name:   "nil config uses defaults",
			config: nil,
			want: &Config{
				HealthCheckInterval: 30 * time.Second,
				ConnectionTimeout:   10 * time.Second,
				MaxRetries:          3,
				RetryBackoff:        5 * time.Second,
			},
		},
		{
			name: "custom config is used",
			config: &Config{
				HealthCheckInterval: 60 * time.Second,
				ConnectionTimeout:   20 * time.Second,
				MaxRetries:          5,
				RetryBackoff:        10 * time.Second,
			},
			want: &Config{
				HealthCheckInterval: 60 * time.Second,
				ConnectionTimeout:   20 * time.Second,
				MaxRetries:          5,
				RetryBackoff:        10 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager(tt.config, nil)
			assert.NotNil(t, m)
			assert.NotNil(t, m.cameras)
			assert.Equal(t, tt.want.HealthCheckInterval, m.config.HealthCheckInterval)
			assert.Equal(t, tt.want.ConnectionTimeout, m.config.ConnectionTimeout)
			assert.Equal(t, tt.want.MaxRetries, m.config.MaxRetries)
			assert.Equal(t, tt.want.RetryBackoff, m.config.RetryBackoff)
		})
	}
}

func TestManager_AddCamera_InvalidInput(t *testing.T) {
	m := NewManager(nil, nil)

	tests := []struct {
		name    string
		camera  *models.Camera
		wantErr bool
	}{
		{
			name:    "nil camera",
			camera:  nil,
			wantErr: true,
		},
		{
			name: "empty ID",
			camera: &models.Camera{
				ID:       "",
				Name:     "Test Camera",
				Host:     "192.168.1.100",
				Port:     80,
				Username: "admin",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "empty host",
			camera: &models.Camera{
				ID:       "cam1",
				Name:     "Test Camera",
				Host:     "",
				Port:     80,
				Username: "admin",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "invalid port",
			camera: &models.Camera{
				ID:       "cam1",
				Name:     "Test Camera",
				Host:     "192.168.1.100",
				Port:     0,
				Username: "admin",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "empty username",
			camera: &models.Camera{
				ID:       "cam1",
				Name:     "Test Camera",
				Host:     "192.168.1.100",
				Port:     80,
				Username: "",
				Password: "password",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := m.AddCamera(ctx, tt.camera)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestManager_RemoveCamera(t *testing.T) {
	m := NewManager(nil, nil)

	// Test removing non-existent camera
	err := m.RemoveCamera("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestManager_GetCamera(t *testing.T) {
	m := NewManager(nil, nil)

	// Test getting non-existent camera
	client, err := m.GetCamera("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, client)
}

func TestManager_ListCameras(t *testing.T) {
	m := NewManager(nil, nil)

	// Test empty list
	cameras := m.ListCameras()
	assert.NotNil(t, cameras)
	assert.Empty(t, cameras)
}

func TestManager_GetCameraStatus(t *testing.T) {
	m := NewManager(nil, nil)

	// Test getting status of non-existent camera
	status, err := m.GetCameraStatus("nonexistent")
	assert.Error(t, err)
	assert.Nil(t, status)
}

func TestCameraClient_CircuitBreaker(t *testing.T) {
	// Create a camera client with circuit breaker open
	client := &CameraClient{
		Camera: &models.Camera{
			ID:   "test-cam",
			Name: "Test Camera",
		},
		CircuitOpen:  true,
		FailureCount: 3,
	}

	ctx := context.Background()

	// Test that operations fail when circuit is open
	t.Run("Reboot fails when circuit open", func(t *testing.T) {
		err := client.Reboot(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})

	t.Run("GetSnapshot fails when circuit open", func(t *testing.T) {
		_, err := client.GetSnapshot(ctx, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})

	t.Run("PTZMove fails when circuit open", func(t *testing.T) {
		err := client.PTZMove(ctx, "Right", 32, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})

	t.Run("PTZStop fails when circuit open", func(t *testing.T) {
		err := client.PTZStop(ctx, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})

	t.Run("PTZGotoPreset fails when circuit open", func(t *testing.T) {
		err := client.PTZGotoPreset(ctx, 0, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})

	t.Run("SetIRLights fails when circuit open", func(t *testing.T) {
		err := client.SetIRLights(ctx, 0, "Auto")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})

	t.Run("SetWhiteLED fails when circuit open", func(t *testing.T) {
		config := &reolink.WhiteLed{
			Channel: 0,
			State:   1,
			Mode:    0,
			Bright:  50,
		}
		err := client.SetWhiteLED(ctx, config)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})

	t.Run("TriggerSiren fails when circuit open", func(t *testing.T) {
		err := client.TriggerSiren(ctx, 0, 3)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})

	t.Run("GetMotionState fails when circuit open", func(t *testing.T) {
		_, err := client.GetMotionState(ctx, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})

	t.Run("GetAIState fails when circuit open", func(t *testing.T) {
		_, err := client.GetAIState(ctx, 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "circuit open")
	})
}

func TestManager_Shutdown(t *testing.T) {
	m := NewManager(nil, nil)

	// Test shutdown with no cameras
	err := m.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestConfig_Defaults(t *testing.T) {
	config := &Config{}

	// Verify zero values
	assert.Equal(t, time.Duration(0), config.HealthCheckInterval)
	assert.Equal(t, time.Duration(0), config.ConnectionTimeout)
	assert.Equal(t, 0, config.MaxRetries)
	assert.Equal(t, time.Duration(0), config.RetryBackoff)
}

func TestCameraStatus(t *testing.T) {
	now := time.Now()
	status := &models.CameraStatus{
		CameraID:    "cam1",
		Status:      "online",
		Model:       "RLC-810A",
		FirmwareVer: "v3.0.0.123",
		Uptime:      3600,
		LastSeen:    now,
	}

	assert.Equal(t, "cam1", status.CameraID)
	assert.Equal(t, "online", status.Status)
	assert.Equal(t, "RLC-810A", status.Model)
	assert.Equal(t, "v3.0.0.123", status.FirmwareVer)
	assert.Equal(t, int64(3600), status.Uptime)
	assert.Equal(t, now, status.LastSeen)
}

// Benchmark tests
func BenchmarkNewManager(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = NewManager(nil, nil)
	}
}

func BenchmarkManager_GetCamera(b *testing.B) {
	m := NewManager(nil, nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = m.GetCamera("nonexistent")
	}
}

func BenchmarkManager_ListCameras(b *testing.B) {
	m := NewManager(nil, nil)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = m.ListCameras()
	}
}
