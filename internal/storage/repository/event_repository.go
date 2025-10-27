package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mosleyit/reolink_server/internal/storage/db"
	"github.com/mosleyit/reolink_server/internal/storage/models"
)

// EventRepository handles event database operations
type EventRepository struct {
	db *db.DB
}

// NewEventRepository creates a new event repository
func NewEventRepository(database *db.DB) *EventRepository {
	return &EventRepository{db: database}
}

// Create creates a new event
func (r *EventRepository) Create(ctx context.Context, event *models.Event) error {
	if event.ID == "" {
		event.ID = uuid.New().String()
	}

	now := time.Now()
	if event.Timestamp.IsZero() {
		event.Timestamp = now
	}
	event.CreatedAt = now

	query := `
		INSERT INTO events (id, camera_id, camera_name, type, severity, timestamp, acknowledged,
			acknowledged_at, metadata, snapshot_path, video_clip_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		event.ID, event.CameraID, event.CameraName, event.Type, event.Severity, event.Timestamp,
		event.Acknowledged, event.AcknowledgedAt, event.Metadata, event.SnapshotPath,
		event.VideoClipURL, event.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

// GetByID retrieves an event by ID
func (r *EventRepository) GetByID(ctx context.Context, id string) (*models.Event, error) {
	query := `
		SELECT id, camera_id, camera_name, type, severity, timestamp, acknowledged, acknowledged_at,
			metadata, snapshot_path, video_clip_url, created_at
		FROM events
		WHERE id = $1
	`

	event := &models.Event{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID, &event.CameraID, &event.CameraName, &event.Type, &event.Severity, &event.Timestamp,
		&event.Acknowledged, &event.AcknowledgedAt, &event.Metadata, &event.SnapshotPath,
		&event.VideoClipURL, &event.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("event not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get event: %w", err)
	}

	return event, nil
}

// ListByCameraID retrieves events for a specific camera
func (r *EventRepository) ListByCameraID(ctx context.Context, cameraID string, limit int, offset int) ([]*models.Event, error) {
	query := `
		SELECT id, camera_id, camera_name, type, severity, timestamp, acknowledged, acknowledged_at,
			metadata, snapshot_path, video_clip_url, created_at
		FROM events
		WHERE camera_id = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, cameraID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

// ListByTimeRange retrieves events within a time range
func (r *EventRepository) ListByTimeRange(ctx context.Context, startTime, endTime time.Time, limit int, offset int) ([]*models.Event, error) {
	query := `
		SELECT id, camera_id, camera_name, type, severity, timestamp, acknowledged, acknowledged_at,
			metadata, snapshot_path, video_clip_url, created_at
		FROM events
		WHERE timestamp >= $1 AND timestamp <= $2
		ORDER BY timestamp DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.QueryContext(ctx, query, startTime, endTime, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list events by time range: %w", err)
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

// ListByType retrieves events by type
func (r *EventRepository) ListByType(ctx context.Context, eventType models.EventType, limit int, offset int) ([]*models.Event, error) {
	query := `
		SELECT id, camera_id, camera_name, type, severity, timestamp, acknowledged, acknowledged_at,
			metadata, snapshot_path, video_clip_url, created_at
		FROM events
		WHERE type = $1
		ORDER BY timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, eventType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list events by type: %w", err)
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

// ListUnacknowledged retrieves unacknowledged events
func (r *EventRepository) ListUnacknowledged(ctx context.Context, limit int, offset int) ([]*models.Event, error) {
	query := `
		SELECT id, camera_id, camera_name, type, severity, timestamp, acknowledged, acknowledged_at,
			metadata, snapshot_path, video_clip_url, created_at
		FROM events
		WHERE acknowledged = FALSE
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list unacknowledged events: %w", err)
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

// List retrieves all events with pagination
func (r *EventRepository) List(ctx context.Context, limit int, offset int) ([]*models.Event, error) {
	query := `
		SELECT id, camera_id, camera_name, type, severity, timestamp, acknowledged, acknowledged_at,
			metadata, snapshot_path, video_clip_url, created_at
		FROM events
		ORDER BY timestamp DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	defer rows.Close()

	return r.scanEvents(rows)
}

// Acknowledge marks an event as acknowledged
func (r *EventRepository) Acknowledge(ctx context.Context, id string) error {
	query := `
		UPDATE events
		SET acknowledged = TRUE, acknowledged_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to acknowledge event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found: %s", id)
	}

	return nil
}

// Delete deletes an event
func (r *EventRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM events WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("event not found: %s", id)
	}

	return nil
}

// DeleteOlderThan deletes events older than the specified time
func (r *EventRepository) DeleteOlderThan(ctx context.Context, olderThan time.Time) (int64, error) {
	query := `DELETE FROM events WHERE timestamp < $1`

	result, err := r.db.ExecContext(ctx, query, olderThan)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old events: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// Count returns the total number of events
func (r *EventRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM events`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}

	return count, nil
}

// CountByCameraID returns the number of events for a specific camera
func (r *EventRepository) CountByCameraID(ctx context.Context, cameraID string) (int, error) {
	query := `SELECT COUNT(*) FROM events WHERE camera_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, cameraID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count events: %w", err)
	}

	return count, nil
}

// scanEvents is a helper function to scan multiple events from rows
func (r *EventRepository) scanEvents(rows *sql.Rows) ([]*models.Event, error) {
	events := []*models.Event{}
	for rows.Next() {
		event := &models.Event{}
		err := rows.Scan(
			&event.ID, &event.CameraID, &event.CameraName, &event.Type, &event.Severity, &event.Timestamp,
			&event.Acknowledged, &event.AcknowledgedAt, &event.Metadata, &event.SnapshotPath,
			&event.VideoClipURL, &event.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating events: %w", err)
	}

	return events, nil
}
