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

// CameraGroupRepository handles camera group database operations
type CameraGroupRepository struct {
	db *db.DB
}

// NewCameraGroupRepository creates a new camera group repository
func NewCameraGroupRepository(database *db.DB) *CameraGroupRepository {
	return &CameraGroupRepository{db: database}
}

// Create creates a new camera group
func (r *CameraGroupRepository) Create(ctx context.Context, group *models.CameraGroup) error {
	if group.ID == "" {
		group.ID = uuid.New().String()
	}

	now := time.Now()
	group.CreatedAt = now
	group.UpdatedAt = now

	query := `
		INSERT INTO camera_groups (id, name, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		group.ID, group.Name, group.Description, group.CreatedAt, group.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create camera group: %w", err)
	}

	return nil
}

// GetByID retrieves a camera group by ID
func (r *CameraGroupRepository) GetByID(ctx context.Context, id string) (*models.CameraGroup, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM camera_groups
		WHERE id = $1
	`

	group := &models.CameraGroup{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("camera group not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get camera group: %w", err)
	}

	return group, nil
}

// GetByName retrieves a camera group by name
func (r *CameraGroupRepository) GetByName(ctx context.Context, name string) (*models.CameraGroup, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM camera_groups
		WHERE name = $1
	`

	group := &models.CameraGroup{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("camera group not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get camera group: %w", err)
	}

	return group, nil
}

// List retrieves all camera groups
func (r *CameraGroupRepository) List(ctx context.Context) ([]*models.CameraGroup, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM camera_groups
		ORDER BY name
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list camera groups: %w", err)
	}
	defer rows.Close()

	groups := []*models.CameraGroup{}
	for rows.Next() {
		group := &models.CameraGroup{}
		err := rows.Scan(
			&group.ID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan camera group: %w", err)
		}
		groups = append(groups, group)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating camera groups: %w", err)
	}

	return groups, nil
}

// Update updates a camera group
func (r *CameraGroupRepository) Update(ctx context.Context, group *models.CameraGroup) error {
	query := `
		UPDATE camera_groups
		SET name = $2, description = $3
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, group.ID, group.Name, group.Description)
	if err != nil {
		return fmt.Errorf("failed to update camera group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("camera group not found: %s", group.ID)
	}

	return nil
}

// Delete deletes a camera group
func (r *CameraGroupRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM camera_groups WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete camera group: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("camera group not found: %s", id)
	}

	return nil
}

// Count returns the total number of camera groups
func (r *CameraGroupRepository) Count(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM camera_groups`

	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count camera groups: %w", err)
	}

	return count, nil
}

