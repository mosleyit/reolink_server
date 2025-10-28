package events

import (
	"context"
	"testing"
	"time"

	"github.com/mosleyit/reolink_server/internal/camera"
	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockSubscriber implements the Subscriber interface for testing
type MockSubscriber struct {
	events []*models.Event
}

func (m *MockSubscriber) OnEvent(event *models.Event) error {
	m.events = append(m.events, event)
	return nil
}

func TestNewProcessor(t *testing.T) {
	t.Run("with nil config uses defaults", func(t *testing.T) {
		manager := camera.NewManager(nil, nil)
		processor := NewProcessor(manager, nil)

		assert.NotNil(t, processor)
		assert.NotNil(t, processor.config)
		assert.Equal(t, 5*time.Second, processor.config.PollInterval)
		assert.Equal(t, 1000, processor.config.EventBufferSize)
	})

	t.Run("with custom config", func(t *testing.T) {
		manager := camera.NewManager(nil, nil)
		config := &Config{
			PollInterval:    10 * time.Second,
			EventBufferSize: 500,
		}
		processor := NewProcessor(manager, config)

		assert.NotNil(t, processor)
		assert.Equal(t, 10*time.Second, processor.config.PollInterval)
		assert.Equal(t, 500, processor.config.EventBufferSize)
	})
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, 5*time.Second, config.PollInterval)
	assert.Equal(t, 5*time.Second, config.MotionCheckPeriod)
	assert.Equal(t, 10*time.Second, config.AICheckPeriod)
	assert.Equal(t, 1000, config.EventBufferSize)
	assert.Equal(t, 10, config.MaxWorkers)
}

func TestProcessor_Subscribe(t *testing.T) {
	manager := camera.NewManager(nil, nil)
	processor := NewProcessor(manager, nil)

	subscriber1 := &MockSubscriber{}
	subscriber2 := &MockSubscriber{}

	processor.Subscribe(subscriber1)
	processor.Subscribe(subscriber2)

	assert.Len(t, processor.subscribers, 2)
}

func TestProcessor_PublishCameraEvent(t *testing.T) {
	manager := camera.NewManager(nil, nil)
	processor := NewProcessor(manager, nil)

	subscriber := &MockSubscriber{}
	processor.Subscribe(subscriber)

	// Start the dispatcher
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	processor.wg.Add(1)
	go processor.dispatchEvents(ctx)

	// Give dispatcher time to start
	time.Sleep(100 * time.Millisecond)

	// Publish an event
	processor.PublishCameraEvent("cam-123", "Test Camera", models.EventCameraOnline)

	// Wait for event to be processed
	time.Sleep(100 * time.Millisecond)

	// Stop processor
	cancel()
	processor.wg.Wait()

	// Verify event was received
	require.Len(t, subscriber.events, 1)
	event := subscriber.events[0]
	assert.Equal(t, "cam-123", event.CameraID)
	assert.Equal(t, "Test Camera", event.CameraName)
	assert.Equal(t, models.EventCameraOnline, event.Type)
	assert.NotEmpty(t, event.ID)
}

func TestProcessor_GetEventChannel(t *testing.T) {
	manager := camera.NewManager(nil, nil)
	processor := NewProcessor(manager, nil)

	eventCh := processor.GetEventChannel()
	assert.NotNil(t, eventCh)
}

func TestProcessor_StartStop(t *testing.T) {
	manager := camera.NewManager(nil, nil)
	processor := NewProcessor(manager, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start processor
	err := processor.Start(ctx)
	assert.NoError(t, err)

	// Give it time to start
	time.Sleep(100 * time.Millisecond)

	// Stop processor
	err = processor.Stop()
	assert.NoError(t, err)
}

func TestProcessor_PublishEvent(t *testing.T) {
	manager := camera.NewManager(nil, nil)
	config := &Config{
		EventBufferSize: 10,
	}
	processor := NewProcessor(manager, config)

	event := &models.Event{
		ID:         "evt-123",
		CameraID:   "cam-123",
		CameraName: "Test Camera",
		Type:       models.EventMotionDetected,
		Timestamp:  time.Now(),
	}

	// Publish event
	processor.publishEvent(event)

	// Read from channel
	select {
	case receivedEvent := <-processor.eventCh:
		assert.Equal(t, event.ID, receivedEvent.ID)
		assert.Equal(t, event.CameraID, receivedEvent.CameraID)
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for event")
	}
}

func TestProcessor_PublishEvent_FullChannel(t *testing.T) {
	manager := camera.NewManager(nil, nil)
	config := &Config{
		EventBufferSize: 2,
	}
	processor := NewProcessor(manager, config)

	// Fill the channel
	for i := 0; i < 2; i++ {
		event := &models.Event{
			ID:       "evt-" + string(rune(i)),
			CameraID: "cam-123",
			Type:     models.EventMotionDetected,
		}
		processor.publishEvent(event)
	}

	// Try to publish one more (should be dropped)
	event := &models.Event{
		ID:       "evt-dropped",
		CameraID: "cam-123",
		Type:     models.EventMotionDetected,
	}
	processor.publishEvent(event)

	// Channel should still have only 2 events
	assert.Len(t, processor.eventCh, 2)
}

func TestProcessor_NotifySubscribers(t *testing.T) {
	manager := camera.NewManager(nil, nil)
	processor := NewProcessor(manager, nil)

	subscriber1 := &MockSubscriber{}
	subscriber2 := &MockSubscriber{}

	processor.Subscribe(subscriber1)
	processor.Subscribe(subscriber2)

	event := &models.Event{
		ID:       "evt-123",
		CameraID: "cam-123",
		Type:     models.EventMotionDetected,
	}

	processor.notifySubscribers(event)

	// Both subscribers should receive the event
	assert.Len(t, subscriber1.events, 1)
	assert.Len(t, subscriber2.events, 1)
	assert.Equal(t, event.ID, subscriber1.events[0].ID)
	assert.Equal(t, event.ID, subscriber2.events[0].ID)
}

func TestConfig_Defaults(t *testing.T) {
	config := DefaultConfig()

	assert.Greater(t, config.PollInterval, time.Duration(0))
	assert.Greater(t, config.MotionCheckPeriod, time.Duration(0))
	assert.Greater(t, config.AICheckPeriod, time.Duration(0))
	assert.Greater(t, config.EventBufferSize, 0)
	assert.Greater(t, config.MaxWorkers, 0)
}

// Benchmark tests
func BenchmarkProcessor_PublishEvent(b *testing.B) {
	manager := camera.NewManager(nil, nil)
	processor := NewProcessor(manager, nil)

	event := &models.Event{
		ID:       "evt-123",
		CameraID: "cam-123",
		Type:     models.EventMotionDetected,
	}

	// Drain events in background
	go func() {
		for range processor.eventCh {
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.publishEvent(event)
	}
}

func BenchmarkProcessor_NotifySubscribers(b *testing.B) {
	manager := camera.NewManager(nil, nil)
	processor := NewProcessor(manager, nil)

	// Add 10 subscribers
	for i := 0; i < 10; i++ {
		processor.Subscribe(&MockSubscriber{})
	}

	event := &models.Event{
		ID:       "evt-123",
		CameraID: "cam-123",
		Type:     models.EventMotionDetected,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.notifySubscribers(event)
	}
}

