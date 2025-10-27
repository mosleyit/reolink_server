package events

import (
	"context"
	"fmt"
	"time"

	"github.com/mosleyit/reolink_server/internal/logger"
	"github.com/mosleyit/reolink_server/internal/storage/models"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Store handles event persistence using Redis Streams
type Store struct {
	client     *redis.Client
	streamName string
}

// StoreConfig holds store configuration
type StoreConfig struct {
	RedisAddr     string
	RedisPassword string
	RedisDB       int
	StreamName    string
}

// DefaultStoreConfig returns default store configuration
func DefaultStoreConfig() *StoreConfig {
	return &StoreConfig{
		RedisAddr:  "localhost:6379",
		RedisDB:    0,
		StreamName: "reolink:events",
	}
}

// NewStore creates a new event store
func NewStore(config *StoreConfig) (*Store, error) {
	if config == nil {
		config = DefaultStoreConfig()
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("Connected to Redis for event storage",
		zap.String("addr", config.RedisAddr),
		zap.String("stream", config.StreamName))

	return &Store{
		client:     client,
		streamName: config.StreamName,
	}, nil
}

// OnEvent implements the Subscriber interface
func (s *Store) OnEvent(event *models.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.SaveEvent(ctx, event)
}

// SaveEvent saves an event to Redis Stream
func (s *Store) SaveEvent(ctx context.Context, event *models.Event) error {
	eventData := map[string]interface{}{
		"id":           event.ID,
		"camera_id":    event.CameraID,
		"camera_name":  event.CameraName,
		"type":         string(event.Type),
		"timestamp":    event.Timestamp.Format(time.RFC3339Nano),
		"acknowledged": fmt.Sprintf("%t", event.Acknowledged),
		"created_at":   event.CreatedAt.Format(time.RFC3339Nano),
	}

	if event.Metadata != "" {
		eventData["metadata"] = event.Metadata
	}

	if event.SnapshotPath != "" {
		eventData["snapshot_path"] = event.SnapshotPath
	}

	// Add to Redis Stream
	_, err := s.client.XAdd(ctx, &redis.XAddArgs{
		Stream: s.streamName,
		Values: eventData,
	}).Result()

	if err != nil {
		logger.Error("Failed to save event to Redis",
			zap.String("event_id", event.ID),
			zap.Error(err))
		return fmt.Errorf("failed to save event: %w", err)
	}

	logger.Debug("Event saved to Redis",
		zap.String("event_id", event.ID),
		zap.String("stream", s.streamName))

	return nil
}

// GetEvents retrieves events from Redis Stream
func (s *Store) GetEvents(ctx context.Context, start, end string, count int64) ([]*models.Event, error) {
	if start == "" {
		start = "-" // Beginning of stream
	}
	if end == "" {
		end = "+" // End of stream
	}

	messages, err := s.client.XRevRange(ctx, s.streamName, end, start).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	events := make([]*models.Event, 0, len(messages))
	for _, msg := range messages {
		event, err := s.messageToEvent(msg)
		if err != nil {
			logger.Warn("Failed to parse event",
				zap.String("message_id", msg.ID),
				zap.Error(err))
			continue
		}
		events = append(events, event)

		if count > 0 && int64(len(events)) >= count {
			break
		}
	}

	return events, nil
}

// GetEventsByCameraID retrieves events for a specific camera
func (s *Store) GetEventsByCameraID(ctx context.Context, cameraID string, count int64) ([]*models.Event, error) {
	// Get all recent events and filter by camera ID
	// Note: For production, consider using a separate stream per camera or secondary indexing
	allEvents, err := s.GetEvents(ctx, "", "", 1000)
	if err != nil {
		return nil, err
	}

	events := make([]*models.Event, 0)
	for _, event := range allEvents {
		if event.CameraID == cameraID {
			events = append(events, event)
			if count > 0 && int64(len(events)) >= count {
				break
			}
		}
	}

	return events, nil
}

// GetEventsByType retrieves events of a specific type
func (s *Store) GetEventsByType(ctx context.Context, eventType models.EventType, count int64) ([]*models.Event, error) {
	allEvents, err := s.GetEvents(ctx, "", "", 1000)
	if err != nil {
		return nil, err
	}

	events := make([]*models.Event, 0)
	for _, event := range allEvents {
		if event.Type == eventType {
			events = append(events, event)
			if count > 0 && int64(len(events)) >= count {
				break
			}
		}
	}

	return events, nil
}

// StreamEvents streams events in real-time using Redis Streams consumer
func (s *Store) StreamEvents(ctx context.Context, lastID string) (<-chan *models.Event, error) {
	eventCh := make(chan *models.Event, 100)

	if lastID == "" {
		lastID = "$" // Start from new messages
	}

	go func() {
		defer close(eventCh)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				// Block for 1 second waiting for new messages
				streams, err := s.client.XRead(ctx, &redis.XReadArgs{
					Streams: []string{s.streamName, lastID},
					Block:   1 * time.Second,
					Count:   10,
				}).Result()

				if err != nil {
					if err != redis.Nil {
						logger.Error("Failed to read from stream", zap.Error(err))
					}
					continue
				}

				for _, stream := range streams {
					for _, msg := range stream.Messages {
						event, err := s.messageToEvent(msg)
						if err != nil {
							logger.Warn("Failed to parse event", zap.Error(err))
							continue
						}

						select {
						case eventCh <- event:
							lastID = msg.ID
						case <-ctx.Done():
							return
						}
					}
				}
			}
		}
	}()

	return eventCh, nil
}

// messageToEvent converts a Redis message to an Event
func (s *Store) messageToEvent(msg redis.XMessage) (*models.Event, error) {
	event := &models.Event{}

	if id, ok := msg.Values["id"].(string); ok {
		event.ID = id
	}

	if cameraID, ok := msg.Values["camera_id"].(string); ok {
		event.CameraID = cameraID
	}

	if cameraName, ok := msg.Values["camera_name"].(string); ok {
		event.CameraName = cameraName
	}

	if eventType, ok := msg.Values["type"].(string); ok {
		event.Type = models.EventType(eventType)
	}

	if timestamp, ok := msg.Values["timestamp"].(string); ok {
		if t, err := time.Parse(time.RFC3339Nano, timestamp); err == nil {
			event.Timestamp = t
		}
	}

	if acknowledged, ok := msg.Values["acknowledged"].(string); ok {
		event.Acknowledged = acknowledged == "true"
	}

	if metadata, ok := msg.Values["metadata"].(string); ok {
		event.Metadata = metadata
	}

	if snapshotPath, ok := msg.Values["snapshot_path"].(string); ok {
		event.SnapshotPath = snapshotPath
	}

	if createdAt, ok := msg.Values["created_at"].(string); ok {
		if t, err := time.Parse(time.RFC3339Nano, createdAt); err == nil {
			event.CreatedAt = t
		}
	}

	return event, nil
}

// TrimStream trims old events from the stream
func (s *Store) TrimStream(ctx context.Context, maxLen int64) error {
	_, err := s.client.XTrimMaxLen(ctx, s.streamName, maxLen).Result()
	if err != nil {
		return fmt.Errorf("failed to trim stream: %w", err)
	}

	logger.Info("Trimmed event stream",
		zap.String("stream", s.streamName),
		zap.Int64("max_len", maxLen))

	return nil
}

// Close closes the Redis connection
func (s *Store) Close() error {
	return s.client.Close()
}
