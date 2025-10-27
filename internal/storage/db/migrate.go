package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"go.uber.org/zap"

	"github.com/mosleyit/reolink_server/internal/logger"
)

// RunMigrations runs all pending database migrations
func (db *DB) RunMigrations(ctx context.Context, migrationsPath string) error {
	logger.Info("Running database migrations", zap.String("path", migrationsPath))

	// Create migrations table if it doesn't exist
	if err := db.createMigrationsTable(ctx); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get list of migration files
	files, err := getMigrationFiles(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get applied migrations
	applied, err := db.getAppliedMigrations(ctx)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Apply pending migrations
	for _, file := range files {
		if _, ok := applied[file]; ok {
			logger.Debug("Migration already applied", zap.String("file", file))
			continue
		}

		logger.Info("Applying migration", zap.String("file", file))

		// Read migration file
		content, err := os.ReadFile(filepath.Join(migrationsPath, file))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		// Execute migration
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		// Record migration
		if err := db.recordMigration(ctx, file); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", file, err)
		}

		logger.Info("Migration applied successfully", zap.String("file", file))
	}

	logger.Info("All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (db *DB) createMigrationsTable(ctx context.Context) error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			version VARCHAR(255) NOT NULL UNIQUE,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		)
	`

	_, err := db.ExecContext(ctx, query)
	return err
}

// getAppliedMigrations returns a map of applied migration versions
func (db *DB) getAppliedMigrations(ctx context.Context) (map[string]bool, error) {
	query := `SELECT version FROM schema_migrations ORDER BY version`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// recordMigration records a migration as applied
func (db *DB) recordMigration(ctx context.Context, version string) error {
	query := `INSERT INTO schema_migrations (version) VALUES ($1)`
	_, err := db.ExecContext(ctx, query, version)
	return err
}

// getMigrationFiles returns a sorted list of .up.sql migration files
func getMigrationFiles(migrationsPath string) ([]string, error) {
	entries, err := os.ReadDir(migrationsPath)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(name, ".up.sql") {
			files = append(files, name)
		}
	}

	sort.Strings(files)
	return files, nil
}

