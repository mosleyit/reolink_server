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

// RecordingRepository handles recording database operations
type RecordingRepository struct {
	db *db.DB
}

// NewRecordingRepository creates a new recording repository
func NewRecordingRepository(database *db.DB) *RecordingRepository {
	return &RecordingRepository{db: database}
}

// Create creates a new recording
func (r *RecordingRepository) Create(ctx context.Context, recording *models.Recording) error {
	if recording.ID == "" {
		recording.ID = uuid.New().String()
	}

	now := time.Now()
	recording.CreatedAt = now

	query := `
		INSERT INTO recordings (id, camera_id, file_name, file_size, start_time, end_time, 
			duration, stream_type, recording_type, storage_path, thumbnail_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		recording.ID, recording.CameraID, recording.FileName, recording.FileSize,
		recording.StartTime, recording.EndTime, recording.Duration, recording.StreamType,
		recording.RecordingType, recording.StoragePath, recording.ThumbnailURL, recording.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create recording: %w", err)
	}

	return nil
}

// GetByID retrieves a recording by ID
func (r *RecordingRepository) GetByID(ctx context.Context, id string) (*models.Recording, error) {
	query := `
		SELECT id, camera_id, file_name, file_size, start_time, end_time, duration,
			stream_type, recording_type, storage_path, thumbnail_url, created_at
		FROM recordings
		WHERE id = $1
	`

	recording := &models.Recording{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&recording.ID, &recording.CameraID, &recording.FileName, &recording.FileSize,
		&recording.StartTime, &recording.EndTime, &recording.Duration, &recording.StreamType,
		&recording.RecordingType, &recording.StoragePath, &recording.ThumbnailURL, &recording.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("recording not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get recording: %w", err)
	}

	return recording, nil
}

// ListByCameraID retrieves recordings for a specific camera
func (r *RecordingRepository) ListByCameraID(ctx context.Context, cameraID string, limit int, offset int) ([]*models.Recording, error) {
	query := `
		SELECT id, camera_id, file_name, file_size, start_time, end_time, duration,
			stream_type, recording_type, storage_path, thumbnail_url, created_at
		FROM recordings
		WHERE camera_id = $1
		ORDER BY start_time DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, cameraID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list recordings: %w", err)
	}
	defer rows.Close()

	return r.scanRecordings(rows)
}

// ListByTimeRange retrieves recordings within a time range
func (r *RecordingRepository) ListByTimeRange(ctx context.Context, cameraID string, startTime, endTime time.Time, limit int, offset int) ([]*models.Recording, error) {
	query := `
		SELECT id, camera_id, file_name, file_size, start_time, end_time, duration,
			stream_type, recording_type, storage_path, thumbnail_url, created_at
		FROM recordings
		WHERE camera_id = $1 AND start_time >= $2 AND end_time <= $3
		ORDER BY start_time DESC
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.QueryContext(ctx, query, cameraID, startTime, endTime, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list recordings by time range: %w", err)
	}
	defer rows.Close()

	return r.scanRecordings(rows)
}

// Search searches recordings with flexible filters
func (r *RecordingRepository) Search(ctx context.Context, req *models.RecordingSearchRequest) ([]*models.Recording, error) {
	query := `
		SELECT id, camera_id, file_name, file_size, start_time, end_time, duration,
			stream_type, recording_type, storage_path, thumbnail_url, created_at
		FROM recordings
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if req.CameraID != nil {
		query += fmt.Sprintf(" AND camera_id = $%d", argCount)
		args = append(args, *req.CameraID)
		argCount++
	}

	if req.StartTime != nil {
		query += fmt.Sprintf(" AND start_time >= $%d", argCount)
		args = append(args, *req.StartTime)
		argCount++
	}

	if req.EndTime != nil {
		query += fmt.Sprintf(" AND end_time <= $%d", argCount)
		args = append(args, *req.EndTime)
		argCount++
	}

	if req.RecordingType != nil {
		query += fmt.Sprintf(" AND recording_type = $%d", argCount)
		args = append(args, *req.RecordingType)
		argCount++
	}

	if req.StreamType != nil {
		query += fmt.Sprintf(" AND stream_type = $%d", argCount)
		args = append(args, *req.StreamType)
		argCount++
	}

	query += " ORDER BY start_time DESC"

	if req.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, req.Limit)
		argCount++
	}

	if req.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, req.Offset)
		argCount++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search recordings: %w", err)
	}
	defer rows.Close()

	return r.scanRecordings(rows)
}

// Delete deletes a recording
func (r *RecordingRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM recordings WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recording: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("recording not found: %s", id)
	}

	return nil
}

// DeleteOlderThan deletes recordings older than the specified time
func (r *RecordingRepository) DeleteOlderThan(ctx context.Context, olderThan time.Time) (int64, error) {
	query := `DELETE FROM recordings WHERE end_time < $1`

	result, err := r.db.ExecContext(ctx, query, olderThan)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old recordings: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// Count returns the total number of recordings
func (r *RecordingRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM recordings`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count recordings: %w", err)
	}

	return count, nil
}

// CountByCameraID returns the number of recordings for a specific camera
func (r *RecordingRepository) CountByCameraID(ctx context.Context, cameraID string) (int, error) {
	query := `SELECT COUNT(*) FROM recordings WHERE camera_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, cameraID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count recordings: %w", err)
	}

	return count, nil
}

// GetTotalSize returns the total size of all recordings in bytes
func (r *RecordingRepository) GetTotalSize(ctx context.Context) (int64, error) {
	query := `SELECT COALESCE(SUM(file_size), 0) FROM recordings`

	var totalSize int64
	err := r.db.QueryRowContext(ctx, query).Scan(&totalSize)
	if err != nil {
		return 0, fmt.Errorf("failed to get total size: %w", err)
	}

	return totalSize, nil
}

// GetTotalSizeByCameraID returns the total size of recordings for a specific camera
func (r *RecordingRepository) GetTotalSizeByCameraID(ctx context.Context, cameraID string) (int64, error) {
	query := `SELECT COALESCE(SUM(file_size), 0) FROM recordings WHERE camera_id = $1`

	var totalSize int64
	err := r.db.QueryRowContext(ctx, query, cameraID).Scan(&totalSize)
	if err != nil {
		return 0, fmt.Errorf("failed to get total size: %w", err)
	}

	return totalSize, nil
}

// scanRecordings is a helper function to scan multiple recordings from rows
func (r *RecordingRepository) scanRecordings(rows *sql.Rows) ([]*models.Recording, error) {
	recordings := []*models.Recording{}
	for rows.Next() {
		recording := &models.Recording{}
		err := rows.Scan(
			&recording.ID, &recording.CameraID, &recording.FileName, &recording.FileSize,
			&recording.StartTime, &recording.EndTime, &recording.Duration, &recording.StreamType,
			&recording.RecordingType, &recording.StoragePath, &recording.ThumbnailURL, &recording.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recording: %w", err)
		}
		recordings = append(recordings, recording)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recordings: %w", err)
	}

	return recordings, nil
}
