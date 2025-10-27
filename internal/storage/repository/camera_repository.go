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

// CameraRepository handles camera database operations
type CameraRepository struct {
	db *db.DB
}

// NewCameraRepository creates a new camera repository
func NewCameraRepository(database *db.DB) *CameraRepository {
	return &CameraRepository{db: database}
}

// Create creates a new camera
func (r *CameraRepository) Create(ctx context.Context, camera *models.Camera) error {
	if camera.ID == "" {
		camera.ID = uuid.New().String()
	}

	now := time.Now()
	camera.CreatedAt = now
	camera.UpdatedAt = now

	query := `
		INSERT INTO cameras (id, name, host, port, username, password, use_https, skip_verify,
			status, model, firmware_version, hardware_version, capabilities, tags, group_id,
			last_seen, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
	`

	_, err := r.db.ExecContext(ctx, query,
		camera.ID, camera.Name, camera.Host, camera.Port, camera.Username, camera.Password,
		camera.UseHTTPS, camera.SkipVerify, camera.Status, camera.Model, camera.FirmwareVer,
		camera.HardwareVer, camera.Capabilities, camera.Tags, camera.GroupID,
		camera.LastSeen, camera.CreatedAt, camera.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create camera: %w", err)
	}

	return nil
}

// GetByID retrieves a camera by ID
func (r *CameraRepository) GetByID(ctx context.Context, id string) (*models.Camera, error) {
	query := `
		SELECT id, name, host, port, username, password, use_https, skip_verify,
			status, model, firmware_version, hardware_version, capabilities, tags, group_id,
			last_seen, created_at, updated_at
		FROM cameras
		WHERE id = $1
	`

	camera := &models.Camera{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&camera.ID, &camera.Name, &camera.Host, &camera.Port, &camera.Username, &camera.Password,
		&camera.UseHTTPS, &camera.SkipVerify, &camera.Status, &camera.Model, &camera.FirmwareVer,
		&camera.HardwareVer, &camera.Capabilities, &camera.Tags, &camera.GroupID,
		&camera.LastSeen, &camera.CreatedAt, &camera.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("camera not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get camera: %w", err)
	}

	return camera, nil
}

// GetByHost retrieves a camera by host and port
func (r *CameraRepository) GetByHost(ctx context.Context, host string, port int) (*models.Camera, error) {
	query := `
		SELECT id, name, host, port, username, password, use_https, skip_verify,
			status, model, firmware_version, hardware_version, capabilities, tags, group_id,
			last_seen, created_at, updated_at
		FROM cameras
		WHERE host = $1 AND port = $2
	`

	camera := &models.Camera{}
	err := r.db.QueryRowContext(ctx, query, host, port).Scan(
		&camera.ID, &camera.Name, &camera.Host, &camera.Port, &camera.Username, &camera.Password,
		&camera.UseHTTPS, &camera.SkipVerify, &camera.Status, &camera.Model, &camera.FirmwareVer,
		&camera.HardwareVer, &camera.Capabilities, &camera.Tags, &camera.GroupID,
		&camera.LastSeen, &camera.CreatedAt, &camera.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("camera not found: %s:%d", host, port)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get camera: %w", err)
	}

	return camera, nil
}

// List retrieves all cameras
func (r *CameraRepository) List(ctx context.Context) ([]*models.Camera, error) {
	query := `
		SELECT id, name, host, port, username, password, use_https, skip_verify,
			status, model, firmware_version, hardware_version, capabilities, tags, group_id,
			last_seen, created_at, updated_at
		FROM cameras
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list cameras: %w", err)
	}
	defer rows.Close()

	cameras := []*models.Camera{}
	for rows.Next() {
		camera := &models.Camera{}
		err := rows.Scan(
			&camera.ID, &camera.Name, &camera.Host, &camera.Port, &camera.Username, &camera.Password,
			&camera.UseHTTPS, &camera.SkipVerify, &camera.Status, &camera.Model, &camera.FirmwareVer,
			&camera.HardwareVer, &camera.Capabilities, &camera.Tags, &camera.GroupID,
			&camera.LastSeen, &camera.CreatedAt, &camera.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan camera: %w", err)
		}
		cameras = append(cameras, camera)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cameras: %w", err)
	}

	return cameras, nil
}

// Update updates a camera
func (r *CameraRepository) Update(ctx context.Context, camera *models.Camera) error {
	query := `
		UPDATE cameras
		SET name = $2, host = $3, port = $4, username = $5, password = $6,
			use_https = $7, skip_verify = $8, status = $9, model = $10,
			firmware_version = $11, hardware_version = $12, capabilities = $13,
			tags = $14, group_id = $15, last_seen = $16
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		camera.ID, camera.Name, camera.Host, camera.Port, camera.Username, camera.Password,
		camera.UseHTTPS, camera.SkipVerify, camera.Status, camera.Model, camera.FirmwareVer,
		camera.HardwareVer, camera.Capabilities, camera.Tags, camera.GroupID, camera.LastSeen)

	if err != nil {
		return fmt.Errorf("failed to update camera: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("camera not found: %s", camera.ID)
	}

	return nil
}

// UpdateStatus updates camera status and last seen time
func (r *CameraRepository) UpdateStatus(ctx context.Context, id string, status string, lastSeen time.Time) error {
	query := `
		UPDATE cameras
		SET status = $2, last_seen = $3
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, status, lastSeen)
	if err != nil {
		return fmt.Errorf("failed to update camera status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("camera not found: %s", id)
	}

	return nil
}

// Delete deletes a camera
func (r *CameraRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM cameras WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete camera: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("camera not found: %s", id)
	}

	return nil
}

// Count returns the total number of cameras
func (r *CameraRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM cameras`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count cameras: %w", err)
	}

	return count, nil
}

// ListByStatus retrieves cameras by status
func (r *CameraRepository) ListByStatus(ctx context.Context, status string) ([]*models.Camera, error) {
	query := `
		SELECT id, name, host, port, username, password, use_https, skip_verify,
			status, model, firmware_version, hardware_version, capabilities, tags, group_id,
			last_seen, created_at, updated_at
		FROM cameras
		WHERE status = $1
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list cameras by status: %w", err)
	}
	defer rows.Close()

	cameras := []*models.Camera{}
	for rows.Next() {
		camera := &models.Camera{}
		err := rows.Scan(
			&camera.ID, &camera.Name, &camera.Host, &camera.Port, &camera.Username, &camera.Password,
			&camera.UseHTTPS, &camera.SkipVerify, &camera.Status, &camera.Model, &camera.FirmwareVer,
			&camera.HardwareVer, &camera.Capabilities, &camera.Tags, &camera.GroupID,
			&camera.LastSeen, &camera.CreatedAt, &camera.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan camera: %w", err)
		}
		cameras = append(cameras, camera)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating cameras: %w", err)
	}

	return cameras, nil
}
