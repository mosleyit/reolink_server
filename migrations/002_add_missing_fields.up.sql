-- Add missing fields to cameras table
ALTER TABLE cameras
    ADD COLUMN IF NOT EXISTS capabilities JSONB DEFAULT '{}',
    ADD COLUMN IF NOT EXISTS tags TEXT[] DEFAULT '{}';

-- Add missing fields to events table
ALTER TABLE events
    ADD COLUMN IF NOT EXISTS severity VARCHAR(20) NOT NULL DEFAULT 'info',
    ADD COLUMN IF NOT EXISTS video_clip_url TEXT,
    ADD COLUMN IF NOT EXISTS acknowledged_at TIMESTAMPTZ;

-- Create index on event severity
CREATE INDEX IF NOT EXISTS idx_events_severity ON events(severity, timestamp DESC);

-- Add check constraint for severity values (skip if already exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'check_event_severity') THEN
        ALTER TABLE events ADD CONSTRAINT check_event_severity
        CHECK (severity IN ('info', 'warning', 'critical'));
    END IF;
END $$;

-- Add check constraint for camera status values (skip if already exists)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'check_camera_status') THEN
        ALTER TABLE cameras ADD CONSTRAINT check_camera_status
        CHECK (status IN ('online', 'offline', 'error'));
    END IF;
END $$;

