-- Remove constraints
ALTER TABLE cameras DROP CONSTRAINT IF EXISTS check_camera_status;
ALTER TABLE events DROP CONSTRAINT IF EXISTS check_event_severity;

-- Remove index
DROP INDEX IF EXISTS idx_events_severity;

-- Remove added columns from events table
ALTER TABLE events
    DROP COLUMN IF EXISTS acknowledged_at,
    DROP COLUMN IF EXISTS video_clip_url,
    DROP COLUMN IF EXISTS severity;

-- Remove added columns from cameras table
ALTER TABLE cameras
    DROP COLUMN IF EXISTS tags,
    DROP COLUMN IF EXISTS capabilities;

