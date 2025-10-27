-- Enable TimescaleDB extension
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Cameras table
CREATE TABLE IF NOT EXISTS cameras (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INTEGER NOT NULL DEFAULT 80,
    username VARCHAR(100) NOT NULL,
    password TEXT NOT NULL,
    use_https BOOLEAN DEFAULT FALSE,
    skip_verify BOOLEAN DEFAULT TRUE,
    status VARCHAR(20) DEFAULT 'offline',
    model VARCHAR(100),
    firmware_version VARCHAR(50),
    hardware_version VARCHAR(50),
    last_seen TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(host, port)
);

CREATE INDEX idx_cameras_status ON cameras(status);
CREATE INDEX idx_cameras_last_seen ON cameras(last_seen);

-- Events table (TimescaleDB hypertable)
CREATE TABLE IF NOT EXISTS events (
    id UUID DEFAULT uuid_generate_v4(),
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    camera_name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    acknowledged BOOLEAN DEFAULT FALSE,
    metadata JSONB DEFAULT '{}',
    snapshot_path TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY (id, timestamp)
);

-- Convert events table to TimescaleDB hypertable
SELECT create_hypertable('events', 'timestamp', if_not_exists => TRUE);

-- Create indexes on events
CREATE INDEX idx_events_camera_id ON events(camera_id, timestamp DESC);
CREATE INDEX idx_events_type ON events(type, timestamp DESC);
CREATE INDEX idx_events_acknowledged ON events(acknowledged, timestamp DESC);

-- Recordings table
CREATE TABLE IF NOT EXISTS recordings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    file_name VARCHAR(500) NOT NULL,
    file_size BIGINT,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    duration INTEGER, -- seconds
    stream_type VARCHAR(20), -- main, sub, ext
    recording_type VARCHAR(50), -- timing, motion, ai_people, ai_vehicle, ai_pet
    storage_path TEXT,
    thumbnail_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(camera_id, file_name)
);

CREATE INDEX idx_recordings_camera_id ON recordings(camera_id, start_time DESC);
CREATE INDEX idx_recordings_time_range ON recordings(start_time, end_time);
CREATE INDEX idx_recordings_type ON recordings(recording_type, start_time DESC);

-- Camera configurations table
CREATE TABLE IF NOT EXISTS camera_configs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    camera_id UUID NOT NULL REFERENCES cameras(id) ON DELETE CASCADE,
    config_type VARCHAR(50) NOT NULL, -- osd, image, isp, encoding, etc.
    config_data JSONB NOT NULL DEFAULT '{}',
    version INTEGER DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(camera_id, config_type)
);

CREATE INDEX idx_camera_configs_camera_id ON camera_configs(camera_id);
CREATE INDEX idx_camera_configs_type ON camera_configs(config_type);

-- Event retention policy (delete events older than 90 days)
CREATE OR REPLACE FUNCTION delete_old_events()
RETURNS void AS $$
BEGIN
    DELETE FROM events WHERE timestamp < NOW() - INTERVAL '90 days';
END;
$$ LANGUAGE plpgsql;

-- Create a scheduled job to run retention policy daily (requires TimescaleDB)
-- This will be executed by the application or a cron job
-- SELECT add_job('delete_old_events', '1 day');

-- Updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply updated_at trigger to cameras
CREATE TRIGGER update_cameras_updated_at
    BEFORE UPDATE ON cameras
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply updated_at trigger to camera_configs
CREATE TRIGGER update_camera_configs_updated_at
    BEFORE UPDATE ON camera_configs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

