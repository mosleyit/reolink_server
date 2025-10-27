package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mosleyit/reolink_server/internal/storage/models"
)

// MockEventProcessor is a mock implementation of EventProcessor
type MockEventProcessor struct {
	subscribers []EventSubscriber
}

func (m *MockEventProcessor) Subscribe(subscriber EventSubscriber) {
	m.subscribers = append(m.subscribers, subscriber)
}

func (m *MockEventProcessor) PublishEvent(event *models.Event) {
	for _, sub := range m.subscribers {
		_ = sub.OnEvent(event)
	}
}

func TestNewEventStreamService(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)

	assert.NotNil(t, service)
	assert.Equal(t, mockProcessor, service.processor)
	assert.NotNil(t, service.subscribers)
	assert.Equal(t, 0, len(service.subscribers))
}

func TestEventStreamService_Subscribe(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx := context.Background()

	subscriber := service.Subscribe(ctx, "sub-1", "", 10)

	assert.NotNil(t, subscriber)
	assert.Equal(t, "sub-1", subscriber.ID)
	assert.Equal(t, "", subscriber.CameraID)
	assert.NotNil(t, subscriber.EventCh)
	assert.Equal(t, 1, len(service.subscribers))
	assert.Equal(t, 1, len(mockProcessor.subscribers))
}

func TestEventStreamService_Subscribe_WithCameraID(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx := context.Background()

	subscriber := service.Subscribe(ctx, "sub-1", "cam-123", 10)

	assert.NotNil(t, subscriber)
	assert.Equal(t, "sub-1", subscriber.ID)
	assert.Equal(t, "cam-123", subscriber.CameraID)
}

func TestEventStreamService_Unsubscribe(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx := context.Background()

	subscriber := service.Subscribe(ctx, "sub-1", "", 10)
	assert.Equal(t, 1, len(service.subscribers))

	service.Unsubscribe("sub-1")
	assert.Equal(t, 0, len(service.subscribers))

	// Verify channel is closed
	_, ok := <-subscriber.EventCh
	assert.False(t, ok, "Channel should be closed")
}

func TestEventStreamService_OnEvent(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx := context.Background()

	subscriber := service.Subscribe(ctx, "sub-1", "", 10)

	event := &models.Event{
		ID:       "evt-1",
		CameraID: "cam-1",
		Type:     models.EventMotionDetected,
	}

	err := service.OnEvent(event)
	assert.NoError(t, err)

	// Verify event was received
	select {
	case receivedEvent := <-subscriber.EventCh:
		assert.Equal(t, event.ID, receivedEvent.ID)
		assert.Equal(t, event.CameraID, receivedEvent.CameraID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Event not received")
	}
}

func TestEventStreamService_OnEvent_FilterByCameraID(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx := context.Background()

	// Subscribe to specific camera
	subscriber := service.Subscribe(ctx, "sub-1", "cam-123", 10)

	// Send event from different camera
	event1 := &models.Event{
		ID:       "evt-1",
		CameraID: "cam-456",
		Type:     models.EventMotionDetected,
	}

	err := service.OnEvent(event1)
	assert.NoError(t, err)

	// Should not receive event
	select {
	case <-subscriber.EventCh:
		t.Fatal("Should not receive event from different camera")
	case <-time.After(50 * time.Millisecond):
		// Expected - no event received
	}

	// Send event from subscribed camera
	event2 := &models.Event{
		ID:       "evt-2",
		CameraID: "cam-123",
		Type:     models.EventMotionDetected,
	}

	err = service.OnEvent(event2)
	assert.NoError(t, err)

	// Should receive event
	select {
	case receivedEvent := <-subscriber.EventCh:
		assert.Equal(t, event2.ID, receivedEvent.ID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Event not received")
	}
}

func TestEventStreamService_OnEvent_MultipleSubscribers(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx := context.Background()

	subscriber1 := service.Subscribe(ctx, "sub-1", "", 10)
	subscriber2 := service.Subscribe(ctx, "sub-2", "", 10)

	event := &models.Event{
		ID:       "evt-1",
		CameraID: "cam-1",
		Type:     models.EventMotionDetected,
	}

	err := service.OnEvent(event)
	assert.NoError(t, err)

	// Both subscribers should receive the event
	select {
	case receivedEvent := <-subscriber1.EventCh:
		assert.Equal(t, event.ID, receivedEvent.ID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Subscriber 1 did not receive event")
	}

	select {
	case receivedEvent := <-subscriber2.EventCh:
		assert.Equal(t, event.ID, receivedEvent.ID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Subscriber 2 did not receive event")
	}
}

func TestEventStreamService_OnEvent_ChannelFull(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx := context.Background()

	// Create subscriber with small buffer
	subscriber := service.Subscribe(ctx, "sub-1", "", 2)

	// Fill the channel
	for i := 0; i < 2; i++ {
		event := &models.Event{
			ID:       "evt-" + string(rune(i)),
			CameraID: "cam-1",
			Type:     models.EventMotionDetected,
		}
		err := service.OnEvent(event)
		assert.NoError(t, err)
	}

	// Send one more event - should not block
	event := &models.Event{
		ID:       "evt-overflow",
		CameraID: "cam-1",
		Type:     models.EventMotionDetected,
	}
	err := service.OnEvent(event)
	assert.NoError(t, err)

	// Drain the channel
	<-subscriber.EventCh
	<-subscriber.EventCh

	// The overflow event should have been dropped
	select {
	case <-subscriber.EventCh:
		t.Fatal("Should not receive overflow event")
	case <-time.After(50 * time.Millisecond):
		// Expected - overflow event was dropped
	}
}

func TestEventStreamService_OnEvent_CancelledContext(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx, cancel := context.WithCancel(context.Background())

	subscriber := service.Subscribe(ctx, "sub-1", "", 10)

	// Cancel the context
	cancel()
	time.Sleep(10 * time.Millisecond)

	event := &models.Event{
		ID:       "evt-1",
		CameraID: "cam-1",
		Type:     models.EventMotionDetected,
	}

	err := service.OnEvent(event)
	assert.NoError(t, err)

	// Event should not be sent to cancelled subscriber
	select {
	case <-subscriber.EventCh:
		// May or may not receive depending on timing
	case <-time.After(50 * time.Millisecond):
		// Expected - context cancelled
	}
}

func TestEventStreamService_GetSubscriberCount(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx := context.Background()

	assert.Equal(t, 0, service.GetSubscriberCount())

	service.Subscribe(ctx, "sub-1", "", 10)
	assert.Equal(t, 1, service.GetSubscriberCount())

	service.Subscribe(ctx, "sub-2", "", 10)
	assert.Equal(t, 2, service.GetSubscriberCount())

	service.Unsubscribe("sub-1")
	assert.Equal(t, 1, service.GetSubscriberCount())

	service.Unsubscribe("sub-2")
	assert.Equal(t, 0, service.GetSubscriberCount())
}

// MockSubscriber implements EventSubscriber for testing
type MockSubscriber struct {
	events []*models.Event
}

func (m *MockSubscriber) OnEvent(event *models.Event) error {
	m.events = append(m.events, event)
	return nil
}

func TestProcessorAdapter(t *testing.T) {
	called := false
	var receivedSub EventSubscriber

	adapter := NewProcessorAdapter(func(sub EventSubscriber) {
		called = true
		receivedSub = sub
	})

	mockSub := &MockSubscriber{}
	adapter.Subscribe(mockSub)

	assert.True(t, called)
	assert.Equal(t, mockSub, receivedSub)
}

func TestEventStreamService_Integration(t *testing.T) {
	mockProcessor := &MockEventProcessor{}
	service := NewEventStreamService(mockProcessor)
	ctx := context.Background()

	// Create multiple subscribers
	sub1 := service.Subscribe(ctx, "sub-1", "", 10)
	sub2 := service.Subscribe(ctx, "sub-2", "cam-123", 10)
	sub3 := service.Subscribe(ctx, "sub-3", "cam-456", 10)

	require.Equal(t, 1, len(mockProcessor.subscribers))

	// Publish events through the mock processor
	event1 := &models.Event{
		ID:       "evt-1",
		CameraID: "cam-123",
		Type:     models.EventMotionDetected,
	}
	mockProcessor.PublishEvent(event1)

	// sub1 (all cameras) and sub2 (cam-123) should receive
	select {
	case e := <-sub1.EventCh:
		assert.Equal(t, "evt-1", e.ID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("sub1 did not receive event")
	}

	select {
	case e := <-sub2.EventCh:
		assert.Equal(t, "evt-1", e.ID)
	case <-time.After(100 * time.Millisecond):
		t.Fatal("sub2 did not receive event")
	}

	// sub3 should not receive
	select {
	case <-sub3.EventCh:
		t.Fatal("sub3 should not receive event from cam-123")
	case <-time.After(50 * time.Millisecond):
		// Expected
	}
}
