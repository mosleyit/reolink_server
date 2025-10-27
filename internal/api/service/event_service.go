package service

import (
	"context"
	"time"

	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/mosleyit/reolink_server/internal/storage/repository"
)

// EventService handles event-related operations
type EventService struct {
	eventRepo *repository.EventRepository
}

// NewEventService creates a new event service
func NewEventService(eventRepo *repository.EventRepository) *EventService {
	return &EventService{
		eventRepo: eventRepo,
	}
}

// GetEvent retrieves an event by ID
func (s *EventService) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	return s.eventRepo.GetByID(ctx, id)
}

// ListEvents retrieves all events with pagination
func (s *EventService) ListEvents(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	return s.eventRepo.List(ctx, limit, offset)
}

// ListEventsByTimeRange retrieves events within a time range
func (s *EventService) ListEventsByTimeRange(ctx context.Context, startTime, endTime time.Time, limit, offset int) ([]*models.Event, error) {
	return s.eventRepo.ListByTimeRange(ctx, startTime, endTime, limit, offset)
}

// ListEventsByType retrieves events by type
func (s *EventService) ListEventsByType(ctx context.Context, eventType models.EventType, limit, offset int) ([]*models.Event, error) {
	return s.eventRepo.ListByType(ctx, eventType, limit, offset)
}

// ListUnacknowledgedEvents retrieves unacknowledged events
func (s *EventService) ListUnacknowledgedEvents(ctx context.Context, limit, offset int) ([]*models.Event, error) {
	return s.eventRepo.ListUnacknowledged(ctx, limit, offset)
}

// AcknowledgeEvent marks an event as acknowledged
func (s *EventService) AcknowledgeEvent(ctx context.Context, id string) error {
	return s.eventRepo.Acknowledge(ctx, id)
}

// CountEvents returns the total number of events
func (s *EventService) CountEvents(ctx context.Context) (int, error) {
	return s.eventRepo.Count(ctx)
}

