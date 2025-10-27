package service

import (
	"context"
	"sync"

	"github.com/mosleyit/reolink_server/internal/storage/models"
)

// EventSubscriber interface for event subscription
type EventSubscriber interface {
	OnEvent(event *models.Event) error
}

// EventProcessor interface for event subscription
type EventProcessor interface {
	Subscribe(subscriber EventSubscriber)
}

// EventStreamService manages event streaming to clients
type EventStreamService struct {
	processor   EventProcessor
	subscribers map[string]*EventStreamSubscriber
	mu          sync.RWMutex
}

// EventStreamSubscriber represents a client subscribed to events
type EventStreamSubscriber struct {
	ID       string
	CameraID string // Empty string means all cameras
	EventCh  chan *models.Event
	ctx      context.Context
	cancel   context.CancelFunc
}

// ProcessorAdapter adapts any processor with a Subscribe method to EventProcessor
type ProcessorAdapter struct {
	subscribe func(EventSubscriber)
}

func (p *ProcessorAdapter) Subscribe(subscriber EventSubscriber) {
	p.subscribe(subscriber)
}

// NewProcessorAdapter creates an adapter for processors
func NewProcessorAdapter(subscribe func(EventSubscriber)) *ProcessorAdapter {
	return &ProcessorAdapter{subscribe: subscribe}
}

// NewEventStreamService creates a new event stream service
func NewEventStreamService(processor EventProcessor) *EventStreamService {
	return &EventStreamService{
		processor:   processor,
		subscribers: make(map[string]*EventStreamSubscriber),
	}
}

// Subscribe creates a new event subscription
func (s *EventStreamService) Subscribe(ctx context.Context, subscriberID string, cameraID string, bufferSize int) *EventStreamSubscriber {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Create subscriber context
	subCtx, cancel := context.WithCancel(ctx)

	subscriber := &EventStreamSubscriber{
		ID:       subscriberID,
		CameraID: cameraID,
		EventCh:  make(chan *models.Event, bufferSize),
		ctx:      subCtx,
		cancel:   cancel,
	}

	s.subscribers[subscriberID] = subscriber

	// Register with event processor if this is the first subscriber
	if len(s.subscribers) == 1 {
		s.processor.Subscribe(s)
	}

	return subscriber
}

// Unsubscribe removes an event subscription
func (s *EventStreamService) Unsubscribe(subscriberID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if subscriber, exists := s.subscribers[subscriberID]; exists {
		subscriber.cancel()
		close(subscriber.EventCh)
		delete(s.subscribers, subscriberID)
	}
}

// OnEvent implements the events.Subscriber interface
func (s *EventStreamService) OnEvent(event *models.Event) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Broadcast to all subscribers
	for _, subscriber := range s.subscribers {
		// Filter by camera ID if specified
		if subscriber.CameraID != "" && subscriber.CameraID != event.CameraID {
			continue
		}

		// Non-blocking send
		select {
		case subscriber.EventCh <- event:
			// Event sent successfully
		case <-subscriber.ctx.Done():
			// Subscriber context cancelled, skip
		default:
			// Channel full, skip this event for this subscriber
			// In production, you might want to log this
		}
	}

	return nil
}

// GetSubscriberCount returns the number of active subscribers
func (s *EventStreamService) GetSubscriberCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.subscribers)
}
